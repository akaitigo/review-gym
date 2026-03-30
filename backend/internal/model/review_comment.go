package model

import "time"

// ReviewComment represents a user-submitted review comment on an exercise.
type ReviewComment struct {
	ID         string    `json:"id"`
	ExerciseID string    `json:"exercise_id"`
	UserID     string    `json:"user_id"`
	FilePath   string    `json:"file_path"`
	LineNumber int       `json:"line_number"`
	Content    string    `json:"content"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate checks that the review comment fields satisfy domain constraints.
func (rc *ReviewComment) Validate() error {
	if rc.ExerciseID == "" {
		return &ValidationError{Field: "exercise_id", Message: "must not be empty"}
	}
	if rc.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "must not be empty"}
	}
	if rc.FilePath == "" {
		return &ValidationError{Field: "file_path", Message: "must not be empty"}
	}
	if rc.LineNumber <= 0 {
		return &ValidationError{Field: "line_number", Message: "must be positive"}
	}
	if rc.Content == "" {
		return &ValidationError{Field: "content", Message: "must not be empty"}
	}
	if len(rc.Content) > 5000 {
		return &ValidationError{Field: "content", Message: "must be at most 5,000 characters"}
	}
	if !rc.Category.Valid() {
		return &ValidationError{Field: "category", Message: "invalid category"}
	}
	return nil
}
