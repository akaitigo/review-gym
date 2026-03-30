package model_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestExerciseValidate(t *testing.T) {
	validExercise := model.Exercise{
		Title:       "Test Exercise",
		Description: "A test exercise for validation",
		Difficulty:  model.DifficultyBeginner,
		Category:    model.CategorySecurity,
		Language:    "Go",
		DiffContent: "--- a/file.go\n+++ b/file.go\n@@ -1 +1 @@\n-old\n+new",
		FilePaths:   []string{"file.go"},
		Metadata:    json.RawMessage(`{}`),
	}

	t.Run("valid exercise passes", func(t *testing.T) {
		e := validExercise
		if err := e.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty title fails", func(t *testing.T) {
		e := validExercise
		e.Title = ""
		err := e.Validate()
		if err == nil {
			t.Error("expected error for empty title")
		}
		assertValidationField(t, err, "title")
	})

	t.Run("title too long fails", func(t *testing.T) {
		e := validExercise
		e.Title = strings.Repeat("a", 201)
		err := e.Validate()
		if err == nil {
			t.Error("expected error for long title")
		}
		assertValidationField(t, err, "title")
	})

	t.Run("empty description fails", func(t *testing.T) {
		e := validExercise
		e.Description = ""
		err := e.Validate()
		if err == nil {
			t.Error("expected error for empty description")
		}
		assertValidationField(t, err, "description")
	})

	t.Run("invalid difficulty fails", func(t *testing.T) {
		e := validExercise
		e.Difficulty = "expert"
		err := e.Validate()
		if err == nil {
			t.Error("expected error for invalid difficulty")
		}
		assertValidationField(t, err, "difficulty")
	})

	t.Run("invalid category fails", func(t *testing.T) {
		e := validExercise
		e.Category = "invalid"
		err := e.Validate()
		if err == nil {
			t.Error("expected error for invalid category")
		}
		assertValidationField(t, err, "category")
	})

	t.Run("invalid category tag fails", func(t *testing.T) {
		e := validExercise
		e.CategoryTags = []model.Category{model.CategorySecurity, "invalid"}
		err := e.Validate()
		if err == nil {
			t.Error("expected error for invalid category tag")
		}
		assertValidationField(t, err, "category_tags")
	})

	t.Run("empty language fails", func(t *testing.T) {
		e := validExercise
		e.Language = ""
		err := e.Validate()
		if err == nil {
			t.Error("expected error for empty language")
		}
		assertValidationField(t, err, "language")
	})

	t.Run("empty diff content fails", func(t *testing.T) {
		e := validExercise
		e.DiffContent = ""
		err := e.Validate()
		if err == nil {
			t.Error("expected error for empty diff content")
		}
		assertValidationField(t, err, "diff_content")
	})

	t.Run("diff content too long fails", func(t *testing.T) {
		e := validExercise
		e.DiffContent = strings.Repeat("x", 100_001)
		err := e.Validate()
		if err == nil {
			t.Error("expected error for too long diff content")
		}
		assertValidationField(t, err, "diff_content")
	})
}

func assertValidationField(t *testing.T, err error, expectedField string) {
	t.Helper()
	ve, ok := err.(*model.ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}
	if ve.Field != expectedField {
		t.Errorf("expected field %q, got %q", expectedField, ve.Field)
	}
}
