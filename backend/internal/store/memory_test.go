package store

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestNewMemoryStore(t *testing.T) {
	ms := NewMemoryStore()

	exercises, err := ms.List(ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Fatal("expected seed exercises to be loaded")
	}
}

func TestMemoryStore_List_NoFilter(t *testing.T) {
	ms := NewMemoryStore()

	exercises, err := ms.List(ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) < 10 {
		t.Errorf("expected at least 10 exercises, got %d", len(exercises))
	}
}

func TestMemoryStore_List_FilterByCategory(t *testing.T) {
	ms := NewMemoryStore()
	cat := model.CategorySecurity
	filter := ExerciseFilter{Category: &cat}

	exercises, err := ms.List(filter)
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
	diff := model.DifficultyBeginner
	filter := ExerciseFilter{Difficulty: &diff}

	exercises, err := ms.List(filter)
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

	exercises, _ := ms.List(ExerciseFilter{})
	if len(exercises) == 0 {
		t.Fatal("no exercises loaded")
	}

	first := exercises[0]
	got, err := ms.GetByID(first.ID)
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

	got, err := ms.GetByID("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestMemoryStore_CreateAndListReviewComments(t *testing.T) {
	ms := NewMemoryStore()

	comment := &model.ReviewComment{
		ExerciseID: "00000001",
		UserID:     "user-1",
		FilePath:   "test.go",
		LineNumber: 5,
		Content:    "This is a test comment",
		Category:   model.CategorySecurity,
	}

	if err := ms.Create(comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID == "" {
		t.Error("expected ID to be set")
	}

	comments, err := ms.ListByExerciseAndUser("00000001", "user-1")
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

	comments, err := ms.ListByExerciseAndUser("nonexistent", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestMemoryStore_ListByExercise_ReferenceReviews(t *testing.T) {
	ms := NewMemoryStore()

	reviews, err := ms.ListByExercise("00000001")
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
