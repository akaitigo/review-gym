package model

import (
	"encoding/json"
	"time"
)

// Score represents the result of scoring a user's review against reference reviews.
type Score struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	ExerciseID      string          `json:"exercise_id"`
	PrecisionScore  float64         `json:"precision_score"`
	RecallScore     float64         `json:"recall_score"`
	OverallScore    float64         `json:"overall_score"`
	CategoryScores  json.RawMessage `json:"category_scores"`
	Feedback        string          `json:"feedback,omitempty"`
	AttemptNumber   int             `json:"attempt_number"`
	DurationSeconds int             `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// CategoryScore represents the score for a specific review category.
type CategoryScore struct {
	Category  Category `json:"category"`
	Score     float64  `json:"score"`
	MaxPoints float64  `json:"max_points"`
	Earned    float64  `json:"earned"`
}

// Validate checks that the score fields satisfy domain constraints.
func (s *Score) Validate() error {
	if s.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "must not be empty"}
	}
	if s.ExerciseID == "" {
		return &ValidationError{Field: "exercise_id", Message: "must not be empty"}
	}
	if s.PrecisionScore < 0 || s.PrecisionScore > 100 {
		return &ValidationError{Field: "precision_score", Message: "must be between 0 and 100"}
	}
	if s.RecallScore < 0 || s.RecallScore > 100 {
		return &ValidationError{Field: "recall_score", Message: "must be between 0 and 100"}
	}
	if s.OverallScore < 0 || s.OverallScore > 100 {
		return &ValidationError{Field: "overall_score", Message: "must be between 0 and 100"}
	}
	if s.AttemptNumber < 1 {
		return &ValidationError{Field: "attempt_number", Message: "must be at least 1"}
	}
	return nil
}
