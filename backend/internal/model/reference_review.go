package model

import "time"

// ReferenceReview represents a curated model review comment for an exercise.
// Each reference review is an expected finding that users should ideally discover.
type ReferenceReview struct {
	ID          string    `json:"id"`
	ExerciseID  string    `json:"exercise_id"`
	FilePath    string    `json:"file_path"`
	LineNumber  int       `json:"line_number"`
	Content     string    `json:"content"`
	Category    Category  `json:"category"`
	Severity    Severity  `json:"severity"`
	Explanation string    `json:"explanation"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate checks that the reference review fields satisfy domain constraints.
func (r *ReferenceReview) Validate() error {
	if r.ExerciseID == "" {
		return &ValidationError{Field: "exercise_id", Message: "must not be empty"}
	}
	if r.FilePath == "" {
		return &ValidationError{Field: "file_path", Message: "must not be empty"}
	}
	if r.LineNumber <= 0 {
		return &ValidationError{Field: "line_number", Message: "must be positive"}
	}
	if r.Content == "" {
		return &ValidationError{Field: "content", Message: "must not be empty"}
	}
	if len(r.Content) > 5000 {
		return &ValidationError{Field: "content", Message: "must be at most 5,000 characters"}
	}
	if !r.Category.Valid() {
		return &ValidationError{Field: "category", Message: "invalid category"}
	}
	if !r.Severity.Valid() {
		return &ValidationError{Field: "severity", Message: "invalid severity"}
	}
	if r.Explanation == "" {
		return &ValidationError{Field: "explanation", Message: "must not be empty"}
	}
	return nil
}
