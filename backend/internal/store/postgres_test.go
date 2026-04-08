package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
)

// testDatabaseURL returns the DATABASE_URL for integration tests.
// Tests are skipped if the environment variable is not set.
func testDatabaseURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping PostgreSQL integration test")
	}
	return url
}

// setupTestDB creates a PostgresStore and sets up the schema for testing.
// It truncates all tables before each test for isolation.
func setupTestDB(t *testing.T) *PostgresStore {
	t.Helper()
	url := testDatabaseURL(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ps, err := NewPostgresStore(ctx, url)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	t.Cleanup(func() {
		if closeErr := ps.Close(); closeErr != nil {
			t.Logf("cleanup close error: %v", closeErr)
		}
	})

	// Truncate all tables for test isolation.
	_, err = ps.db.Exec(`
		TRUNCATE TABLE scores, review_comments, reference_reviews, exercises, user_profiles CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	return ps
}

// seedTestExercise inserts a published exercise and returns its ID.
func seedTestExercise(t *testing.T, db *sql.DB, title string) string {
	t.Helper()
	var id string
	err := db.QueryRow(`
		INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
		VALUES ($1, 'Test description', 'beginner', 'security', '[]', 'go', 'diff content', '["test.go"]', '{}', true)
		RETURNING id
	`, title).Scan(&id)
	if err != nil {
		t.Fatalf("failed to seed exercise: %v", err)
	}
	return id
}

// seedTestReferenceReview inserts a reference review for the given exercise.
func seedTestReferenceReview(t *testing.T, db *sql.DB, exerciseID string) string {
	t.Helper()
	var id string
	err := db.QueryRow(`
		INSERT INTO reference_reviews (exercise_id, file_path, line_number, content, category, severity, explanation)
		VALUES ($1, 'test.go', 10, 'SQL injection vulnerability', 'security', 'critical', 'User input concatenated into SQL query')
		RETURNING id
	`, exerciseID).Scan(&id)
	if err != nil {
		t.Fatalf("failed to seed reference review: %v", err)
	}
	return id
}

// seedTestUser inserts a user profile and returns the user ID.
func seedTestUser(t *testing.T, db *sql.DB, name string) string {
	t.Helper()
	var id string
	err := db.QueryRow(`
		INSERT INTO user_profiles (display_name, weakness_categories)
		VALUES ($1, '[]')
		RETURNING id
	`, name).Scan(&id)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return id
}

func TestPostgresStore_List_NoFilter(t *testing.T) {
	ps := setupTestDB(t)

	seedTestExercise(t, ps.db, "Exercise 1")
	seedTestExercise(t, ps.db, "Exercise 2")

	exercises, err := ps.List(ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) != 2 {
		t.Errorf("expected 2 exercises, got %d", len(exercises))
	}
}

func TestPostgresStore_List_FilterByCategory(t *testing.T) {
	ps := setupTestDB(t)

	seedTestExercise(t, ps.db, "Security Exercise")

	// Insert a non-security exercise.
	_, err := ps.db.Exec(`
		INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
		VALUES ('Perf Exercise', 'desc', 'beginner', 'performance', '[]', 'go', 'diff', '[]', '{}', true)
	`)
	if err != nil {
		t.Fatalf("failed to insert exercise: %v", err)
	}

	cat := model.CategorySecurity
	exercises, err := ps.List(ExerciseFilter{Category: &cat})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) != 1 {
		t.Errorf("expected 1 security exercise, got %d", len(exercises))
	}
}

func TestPostgresStore_List_FilterByCategoryTags(t *testing.T) {
	ps := setupTestDB(t)

	// Insert an exercise with category "design" but category_tags includes "security".
	_, err := ps.db.Exec(`
		INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
		VALUES ('Multi-category', 'desc', 'beginner', 'design', '["security"]', 'go', 'diff', '[]', '{}', true)
	`)
	if err != nil {
		t.Fatalf("failed to insert exercise: %v", err)
	}

	cat := model.CategorySecurity
	exercises, err := ps.List(ExerciseFilter{Category: &cat})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) != 1 {
		t.Errorf("expected 1 exercise via category_tags, got %d", len(exercises))
	}
}

func TestPostgresStore_List_FilterByDifficulty(t *testing.T) {
	ps := setupTestDB(t)

	seedTestExercise(t, ps.db, "Beginner Exercise")

	_, err := ps.db.Exec(`
		INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
		VALUES ('Advanced Exercise', 'desc', 'advanced', 'security', '[]', 'go', 'diff', '[]', '{}', true)
	`)
	if err != nil {
		t.Fatalf("failed to insert exercise: %v", err)
	}

	diff := model.DifficultyBeginner
	exercises, err := ps.List(ExerciseFilter{Difficulty: &diff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) != 1 {
		t.Errorf("expected 1 beginner exercise, got %d", len(exercises))
	}
}

func TestPostgresStore_List_HidesUnpublished(t *testing.T) {
	ps := setupTestDB(t)

	seedTestExercise(t, ps.db, "Published Exercise")

	_, err := ps.db.Exec(`
		INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
		VALUES ('Draft Exercise', 'desc', 'beginner', 'security', '[]', 'go', 'diff', '[]', '{}', false)
	`)
	if err != nil {
		t.Fatalf("failed to insert draft exercise: %v", err)
	}

	exercises, err := ps.List(ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) != 1 {
		t.Errorf("expected 1 published exercise, got %d", len(exercises))
	}
}

func TestPostgresStore_GetByID(t *testing.T) {
	ps := setupTestDB(t)

	id := seedTestExercise(t, ps.db, "Test Exercise")

	ex, err := ps.GetByID(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected exercise, got nil")
	}
	if ex.Title != "Test Exercise" {
		t.Errorf("got title %q, want %q", ex.Title, "Test Exercise")
	}
}

func TestPostgresStore_GetByID_NotFound(t *testing.T) {
	ps := setupTestDB(t)

	ex, err := ps.GetByID("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex != nil {
		t.Errorf("expected nil, got %+v", ex)
	}
}

func TestPostgresStore_CreateAndListReviewComments(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID := seedTestExercise(t, ps.db, "Comment Exercise")
	userID := seedTestUser(t, ps.db, "Test User")

	comment := &model.ReviewComment{
		ExerciseID: exerciseID,
		UserID:     userID,
		FilePath:   "test.go",
		LineNumber: 5,
		Content:    "This is a test comment",
		Category:   model.CategorySecurity,
	}

	if err := ps.Create(comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID == "" {
		t.Error("expected ID to be set")
	}

	comments, err := ps.ListByExerciseAndUser(exerciseID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].Content != "This is a test comment" {
		t.Errorf("got content %q, want %q", comments[0].Content, "This is a test comment")
	}
}

func TestPostgresStore_ListByExerciseAndUser_Empty(t *testing.T) {
	ps := setupTestDB(t)

	comments, err := ps.ListByExerciseAndUser("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestPostgresStore_ListByExercise_ReferenceReviews(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID := seedTestExercise(t, ps.db, "Ref Review Exercise")
	seedTestReferenceReview(t, ps.db, exerciseID)

	reviews, err := ps.ListByExercise(exerciseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 reference review, got %d", len(reviews))
	}
	if reviews[0].ExerciseID != exerciseID {
		t.Errorf("got exercise_id %q, want %q", reviews[0].ExerciseID, exerciseID)
	}
}

func TestPostgresStore_SaveAndGetScores(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID := seedTestExercise(t, ps.db, "Score Exercise")
	userID := seedTestUser(t, ps.db, "Score User")

	score := &model.Score{
		UserID:         userID,
		ExerciseID:     exerciseID,
		PrecisionScore: 75.0,
		RecallScore:    50.0,
		OverallScore:   60.0,
		CategoryScores: json.RawMessage(`{}`),
	}

	if err := ps.SaveScore(score); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.ID == "" {
		t.Error("expected ID to be set")
	}
	if score.AttemptNumber != 1 {
		t.Errorf("expected attempt_number 1, got %d", score.AttemptNumber)
	}

	scores, err := ps.GetScoresByExerciseAndUser(exerciseID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].PrecisionScore != 75 {
		t.Errorf("precision = %.1f, want 75", scores[0].PrecisionScore)
	}
}

func TestPostgresStore_SaveScore_AutoIncrementsAttempt(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID := seedTestExercise(t, ps.db, "Multi Score Exercise")
	userID := seedTestUser(t, ps.db, "Multi Score User")

	for i := 1; i <= 3; i++ {
		score := &model.Score{
			UserID:         userID,
			ExerciseID:     exerciseID,
			PrecisionScore: float64(i * 25),
			RecallScore:    float64(i * 20),
			OverallScore:   float64(i * 22),
			CategoryScores: json.RawMessage(`{}`),
		}
		if err := ps.SaveScore(score); err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", i, err)
		}
		if score.AttemptNumber != i {
			t.Errorf("attempt %d: got attempt_number %d", i, score.AttemptNumber)
		}
	}
}

func TestPostgresStore_GetScoresByUser(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID1 := seedTestExercise(t, ps.db, "Exercise A")
	exerciseID2 := seedTestExercise(t, ps.db, "Exercise B")
	userID := seedTestUser(t, ps.db, "User A")
	otherUserID := seedTestUser(t, ps.db, "User B")

	for _, eid := range []string{exerciseID1, exerciseID2} {
		if err := ps.SaveScore(&model.Score{
			UserID:         userID,
			ExerciseID:     eid,
			PrecisionScore: 50,
			RecallScore:    50,
			OverallScore:   50,
			CategoryScores: json.RawMessage(`{}`),
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if err := ps.SaveScore(&model.Score{
		UserID:         otherUserID,
		ExerciseID:     exerciseID1,
		PrecisionScore: 50,
		RecallScore:    50,
		OverallScore:   50,
		CategoryScores: json.RawMessage(`{}`),
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scores, err := ps.GetScoresByUser(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("expected 2 scores, got %d", len(scores))
	}
}

func TestPostgresStore_CountCompletedExercises(t *testing.T) {
	ps := setupTestDB(t)

	exerciseID1 := seedTestExercise(t, ps.db, "Count Exercise 1")
	exerciseID2 := seedTestExercise(t, ps.db, "Count Exercise 2")
	userID := seedTestUser(t, ps.db, "Count User")

	// Score exercise 1 twice and exercise 2 once.
	for _, eid := range []string{exerciseID1, exerciseID1, exerciseID2} {
		if err := ps.SaveScore(&model.Score{
			UserID:         userID,
			ExerciseID:     eid,
			PrecisionScore: 50,
			RecallScore:    50,
			OverallScore:   50,
			CategoryScores: json.RawMessage(`{}`),
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	count, err := ps.CountCompletedExercises(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 completed exercises, got %d", count)
	}
}
