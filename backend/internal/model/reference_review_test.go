package model_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestReferenceReviewValidate(t *testing.T) {
	valid := model.ReferenceReview{
		ExerciseID:  "ex-001",
		FilePath:    "internal/handler/user.go",
		LineNumber:  10,
		Content:     "SQL injection detected",
		Category:    model.CategorySecurity,
		Severity:    model.SeverityCritical,
		Explanation: "Use parameterized queries instead of string concatenation.",
	}

	t.Run("valid review passes", func(t *testing.T) {
		r := valid
		if err := r.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty exercise_id fails", func(t *testing.T) {
		r := valid
		r.ExerciseID = ""
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "exercise_id")
	})

	t.Run("empty file_path fails", func(t *testing.T) {
		r := valid
		r.FilePath = ""
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "file_path")
	})

	t.Run("zero line_number fails", func(t *testing.T) {
		r := valid
		r.LineNumber = 0
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "line_number")
	})

	t.Run("negative line_number fails", func(t *testing.T) {
		r := valid
		r.LineNumber = -1
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "line_number")
	})

	t.Run("empty content fails", func(t *testing.T) {
		r := valid
		r.Content = ""
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "content")
	})

	t.Run("invalid category fails", func(t *testing.T) {
		r := valid
		r.Category = "invalid"
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "category")
	})

	t.Run("invalid severity fails", func(t *testing.T) {
		r := valid
		r.Severity = "invalid"
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "severity")
	})

	t.Run("empty explanation fails", func(t *testing.T) {
		r := valid
		r.Explanation = ""
		err := r.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "explanation")
	})
}
