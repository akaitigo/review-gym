package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
)

// testRedisURL returns the REDIS_URL for integration tests.
// Tests are skipped if the environment variable is not set.
func testRedisURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("TEST_REDIS_URL not set; skipping Redis integration test")
	}
	return url
}

func TestRedisCache_List_CachesResult(t *testing.T) {
	redisURL := testRedisURL(t)

	// Use a spy memory store to track calls.
	spy := &spyExerciseStore{
		exercises: []model.Exercise{
			{
				ID: "1", Title: "Cached Exercise", Description: "desc",
				Difficulty: model.DifficultyBeginner, Category: model.CategorySecurity,
				Language: "go", DiffContent: "diff", IsPublished: true,
				CategoryTags: []model.Category{},
				FilePaths:    []string{},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ms := NewMemoryStore()
	rc, err := NewRedisCache(ctx, redisURL, spy, ms, ms, ms)
	if err != nil {
		t.Fatalf("failed to create redis cache: %v", err)
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			t.Logf("redis close error: %v", closeErr)
		}
	}()

	// First call should hit the spy.
	exercises, err := rc.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("first List call failed: %v", err)
	}
	if len(exercises) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(exercises))
	}
	if spy.listCalls != 1 {
		t.Errorf("expected 1 spy call, got %d", spy.listCalls)
	}

	// Second call should hit cache.
	exercises, err = rc.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("second List call failed: %v", err)
	}
	if len(exercises) != 1 {
		t.Fatalf("expected 1 exercise from cache, got %d", len(exercises))
	}
	if spy.listCalls != 1 {
		t.Errorf("expected spy to still have 1 call (cached), got %d", spy.listCalls)
	}
}

func TestRedisCache_GetByID_CachesResult(t *testing.T) {
	redisURL := testRedisURL(t)

	spy := &spyExerciseStore{
		exercises: []model.Exercise{
			{
				ID: "test-id", Title: "Cached Detail", Description: "desc",
				Difficulty: model.DifficultyBeginner, Category: model.CategorySecurity,
				Language: "go", DiffContent: "diff", IsPublished: true,
				CategoryTags: []model.Category{},
				FilePaths:    []string{},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ms := NewMemoryStore()
	rc, err := NewRedisCache(ctx, redisURL, spy, ms, ms, ms)
	if err != nil {
		t.Fatalf("failed to create redis cache: %v", err)
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			t.Logf("redis close error: %v", closeErr)
		}
	}()

	// First call hits spy.
	ex, err := rc.GetByID(ctx, "test-id")
	if err != nil {
		t.Fatalf("first GetByID call failed: %v", err)
	}
	if ex == nil {
		t.Fatal("expected exercise, got nil")
	}
	if spy.getByIDCalls != 1 {
		t.Errorf("expected 1 spy call, got %d", spy.getByIDCalls)
	}

	// Second call hits cache.
	ex, err = rc.GetByID(ctx, "test-id")
	if err != nil {
		t.Fatalf("second GetByID call failed: %v", err)
	}
	if ex == nil {
		t.Fatal("expected exercise from cache, got nil")
	}
	if spy.getByIDCalls != 1 {
		t.Errorf("expected spy to still have 1 call (cached), got %d", spy.getByIDCalls)
	}
}

func TestRedisCache_GetByID_NotFound(t *testing.T) {
	redisURL := testRedisURL(t)

	spy := &spyExerciseStore{
		exercises: []model.Exercise{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ms := NewMemoryStore()
	rc, err := NewRedisCache(ctx, redisURL, spy, ms, ms, ms)
	if err != nil {
		t.Fatalf("failed to create redis cache: %v", err)
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			t.Logf("redis close error: %v", closeErr)
		}
	}()

	ex, err := rc.GetByID(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex != nil {
		t.Errorf("expected nil, got %+v", ex)
	}
}

func TestRedisCache_DelegatesWriteOperations(t *testing.T) {
	redisURL := testRedisURL(t)

	ms := NewMemoryStore()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rc, err := NewRedisCache(ctx, redisURL, ms, ms, ms, ms)
	if err != nil {
		t.Fatalf("failed to create redis cache: %v", err)
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			t.Logf("redis close error: %v", closeErr)
		}
	}()

	// Create a comment through Redis cache.
	comment := &model.ReviewComment{
		ExerciseID: "00000001",
		UserID:     "user-1",
		FilePath:   "test.go",
		LineNumber: 5,
		Content:    "Test via Redis cache",
		Category:   model.CategorySecurity,
	}
	if err := rc.Create(ctx, comment); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Verify it was written to the underlying store.
	comments, err := rc.ListByExerciseAndUser(ctx, "00000001", "user-1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}

	// Save a score through Redis cache.
	score := &model.Score{
		UserID:         "user-1",
		ExerciseID:     "00000001",
		PrecisionScore: 80,
		RecallScore:    60,
		OverallScore:   68,
	}
	if err := rc.SaveScore(ctx, score); err != nil {
		t.Fatalf("save score failed: %v", err)
	}

	scores, err := rc.GetScoresByExerciseAndUser(ctx, "00000001", "user-1")
	if err != nil {
		t.Fatalf("get scores failed: %v", err)
	}
	if len(scores) != 1 {
		t.Errorf("expected 1 score, got %d", len(scores))
	}
}

// spyExerciseStore tracks calls to List and GetByID for verifying caching behavior.
type spyExerciseStore struct {
	exercises    []model.Exercise
	listCalls    int
	getByIDCalls int
}

func (s *spyExerciseStore) List(_ context.Context, _ ExerciseFilter) ([]model.Exercise, error) {
	s.listCalls++
	return s.exercises, nil
}

func (s *spyExerciseStore) GetByID(_ context.Context, id string) (*model.Exercise, error) {
	s.getByIDCalls++
	for i := range s.exercises {
		if s.exercises[i].ID == id {
			return &s.exercises[i], nil
		}
	}
	return nil, nil
}
