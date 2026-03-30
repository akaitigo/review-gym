package seed_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/seed"
)

func TestAllReturnsAtLeast10Exercises(t *testing.T) {
	exercises := seed.All()
	if len(exercises) < 10 {
		t.Errorf("All() returned %d exercises, want at least 10", len(exercises))
	}
}

func TestAllExercisesHaveAtLeast3Reviews(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		if len(ew.Reviews) < 3 {
			t.Errorf("exercise %q has %d reference reviews, want at least 3",
				ew.Exercise.Title, len(ew.Reviews))
		}
	}
}

func TestAllExercisesArePublished(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		if !ew.Exercise.IsPublished {
			t.Errorf("exercise %q is not published", ew.Exercise.Title)
		}
	}
}

func TestAllExercisesPassValidation(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		e := ew.Exercise
		if err := e.Validate(); err != nil {
			t.Errorf("exercise %q failed validation: %v", e.Title, err)
		}
	}
}

func TestAllReferenceReviewsPassValidation(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		for i, r := range ew.Reviews {
			// Seed reviews don't have ExerciseID set (it's set at insert time)
			// so we set a dummy value for validation purposes
			r.ExerciseID = "dummy-id"
			if err := r.Validate(); err != nil {
				t.Errorf("exercise %q, review[%d] failed validation: %v",
					ew.Exercise.Title, i, err)
			}
		}
	}
}

func TestAllExercisesHaveUniqueTitle(t *testing.T) {
	exercises := seed.All()
	seen := make(map[string]bool)
	for _, ew := range exercises {
		if seen[ew.Exercise.Title] {
			t.Errorf("duplicate exercise title: %q", ew.Exercise.Title)
		}
		seen[ew.Exercise.Title] = true
	}
}

func TestCategoryDistribution(t *testing.T) {
	exercises := seed.All()
	categoryCounts := make(map[string]int)
	for _, ew := range exercises {
		categoryCounts[string(ew.Exercise.Category)]++
	}

	// Ensure at least 2 categories are represented
	if len(categoryCounts) < 2 {
		t.Errorf("only %d categories represented, want at least 2", len(categoryCounts))
	}
}

func TestDifficultyDistribution(t *testing.T) {
	exercises := seed.All()
	difficultyCounts := make(map[string]int)
	for _, ew := range exercises {
		difficultyCounts[string(ew.Exercise.Difficulty)]++
	}

	// Ensure all 3 difficulty levels are represented
	for _, d := range []string{"beginner", "intermediate", "advanced"} {
		if difficultyCounts[d] == 0 {
			t.Errorf("no exercises with difficulty %q", d)
		}
	}
}

func TestReferenceReviewsHaveCategoriesSet(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		for i, r := range ew.Reviews {
			if !r.Category.Valid() {
				t.Errorf("exercise %q, review[%d] has invalid category %q",
					ew.Exercise.Title, i, r.Category)
			}
		}
	}
}

func TestReferenceReviewsHaveSeveritiesSet(t *testing.T) {
	exercises := seed.All()
	for _, ew := range exercises {
		for i, r := range ew.Reviews {
			if !r.Severity.Valid() {
				t.Errorf("exercise %q, review[%d] has invalid severity %q",
					ew.Exercise.Title, i, r.Severity)
			}
		}
	}
}
