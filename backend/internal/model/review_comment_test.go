package model_test

import (
	"strings"
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestReviewCommentValidate(t *testing.T) {
	valid := model.ReviewComment{
		ExerciseID: "ex-001",
		UserID:     "user-001",
		FilePath:   "internal/handler/user.go",
		LineNumber: 10,
		Content:    "This looks like an SQL injection.",
		Category:   model.CategorySecurity,
	}

	t.Run("valid comment passes", func(t *testing.T) {
		rc := valid
		if err := rc.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty exercise_id fails", func(t *testing.T) {
		rc := valid
		rc.ExerciseID = ""
		err := rc.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "exercise_id")
	})

	t.Run("empty user_id fails", func(t *testing.T) {
		rc := valid
		rc.UserID = ""
		err := rc.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "user_id")
	})

	t.Run("content too long fails", func(t *testing.T) {
		rc := valid
		rc.Content = strings.Repeat("x", 5001)
		err := rc.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "content")
	})
}
