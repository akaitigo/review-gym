package store

import (
	"context"
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestNewMemoryStore(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	exercises, err := ms.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Fatal("expected seed exercises to be loaded")
	}
}

func TestMemoryStore_List_NoFilter(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	exercises, err := ms.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) < 10 {
		t.Errorf("expected at least 10 exercises, got %d", len(exercises))
	}
}

func TestMemoryStore_List_FilterByCategory(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()
	cat := model.CategorySecurity
	filter := ExerciseFilter{Category: &cat}

	exercises, err := ms.List(ctx, filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Fatal("expected at least one security exercise")
	}
	for _, ex := range exercises {
		if ex.Category != model.CategorySecurity {
			// Check category_tags
			found := false
			for _, tag := range ex.CategoryTags {
				if tag == model.CategorySecurity {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("exercise %q has category %q but does not match security", ex.Title, ex.Category)
			}
		}
	}
}

func TestMemoryStore_List_FilterByDifficulty(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()
	diff := model.DifficultyBeginner
	filter := ExerciseFilter{Difficulty: &diff}

	exercises, err := ms.List(ctx, filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, ex := range exercises {
		if ex.Difficulty != model.DifficultyBeginner {
			t.Errorf("exercise %q has difficulty %q, want beginner", ex.Title, ex.Difficulty)
		}
	}
}

func TestMemoryStore_GetByID(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	exercises, _ := ms.List(ctx, ExerciseFilter{})
	if len(exercises) == 0 {
		t.Fatal("no exercises loaded")
	}

	first := exercises[0]
	got, err := ms.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected exercise, got nil")
	}
	if got.ID != first.ID {
		t.Errorf("got ID %q, want %q", got.ID, first.ID)
	}
	if got.DiffContent == "" {
		t.Error("expected diff_content to be populated")
	}
}

func TestMemoryStore_GetByID_NotFound(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	got, err := ms.GetByID(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestMemoryStore_CreateAndListReviewComments(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	comment := &model.ReviewComment{
		ExerciseID: "00000001",
		UserID:     "user-1",
		FilePath:   "test.go",
		LineNumber: 5,
		Content:    "This is a test comment",
		Category:   model.CategorySecurity,
	}

	if err := ms.Create(ctx, comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID == "" {
		t.Error("expected ID to be set")
	}

	comments, err := ms.ListByExerciseAndUser(ctx, "00000001", "user-1")
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

func TestMemoryStore_ListByExerciseAndUser_Empty(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	comments, err := ms.ListByExerciseAndUser(ctx, "nonexistent", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestMemoryStore_ListByExercise_ReferenceReviews(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	reviews, err := ms.ListByExercise(ctx, "00000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviews) == 0 {
		t.Fatal("expected reference reviews for first exercise")
	}
	for _, r := range reviews {
		if r.ExerciseID != "00000001" {
			t.Errorf("reference review has exercise_id %q, want %q", r.ExerciseID, "00000001")
		}
	}
}

func TestMemoryStore_SaveAndGetScores(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	score := &model.Score{
		UserID:         "user-1",
		ExerciseID:     "00000001",
		PrecisionScore: 75.0,
		RecallScore:    50.0,
		OverallScore:   60.0,
		AttemptNumber:  1,
	}

	if err := ms.SaveScore(ctx, score); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.ID == "" {
		t.Error("expected ID to be set")
	}

	scores, err := ms.GetScoresByExerciseAndUser(ctx, "00000001", "user-1")
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

func TestMemoryStore_GetScoresByUser(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	score1 := &model.Score{
		UserID:        "user-1",
		ExerciseID:    "00000001",
		AttemptNumber: 1,
	}
	score2 := &model.Score{
		UserID:        "user-1",
		ExerciseID:    "00000002",
		AttemptNumber: 1,
	}
	score3 := &model.Score{
		UserID:        "user-2",
		ExerciseID:    "00000001",
		AttemptNumber: 1,
	}

	_ = ms.SaveScore(ctx, score1)
	_ = ms.SaveScore(ctx, score2)
	_ = ms.SaveScore(ctx, score3)

	scores, err := ms.GetScoresByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("expected 2 scores, got %d", len(scores))
	}
}

func TestMemoryStore_CountCompletedExercises(t *testing.T) {
	ms := NewMemoryStore()
	ctx := context.Background()

	_ = ms.SaveScore(ctx, &model.Score{UserID: "user-1", ExerciseID: "ex-1", AttemptNumber: 1})
	_ = ms.SaveScore(ctx, &model.Score{UserID: "user-1", ExerciseID: "ex-1", AttemptNumber: 2})
	_ = ms.SaveScore(ctx, &model.Score{UserID: "user-1", ExerciseID: "ex-2", AttemptNumber: 1})

	count, err := ms.CountCompletedExercises(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 completed exercises, got %d", count)
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{1, "00000001"},
		{12, "00000012"},
		{100, "00000100"},
	}
	for _, tt := range tests {
		got := generateID(tt.input)
		if got != tt.want {
			t.Errorf("generateID(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
