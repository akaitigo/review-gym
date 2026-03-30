package model_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestScoreValidate(t *testing.T) {
	valid := model.Score{
		UserID:         "user-001",
		ExerciseID:     "ex-001",
		PrecisionScore: 85.5,
		RecallScore:    72.0,
		OverallScore:   78.75,
		AttemptNumber:  1,
	}

	t.Run("valid score passes", func(t *testing.T) {
		s := valid
		if err := s.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty user_id fails", func(t *testing.T) {
		s := valid
		s.UserID = ""
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "user_id")
	})

	t.Run("empty exercise_id fails", func(t *testing.T) {
		s := valid
		s.ExerciseID = ""
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "exercise_id")
	})

	t.Run("precision_score out of range fails", func(t *testing.T) {
		s := valid
		s.PrecisionScore = 101
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "precision_score")
	})

	t.Run("negative recall_score fails", func(t *testing.T) {
		s := valid
		s.RecallScore = -1
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "recall_score")
	})

	t.Run("overall_score over 100 fails", func(t *testing.T) {
		s := valid
		s.OverallScore = 100.01
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "overall_score")
	})

	t.Run("attempt_number zero fails", func(t *testing.T) {
		s := valid
		s.AttemptNumber = 0
		err := s.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "attempt_number")
	})
}
