package model

import (
	"encoding/json"
	"time"
)

// Exercise represents a PR review exercise with anonymized diff data.
type Exercise struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Difficulty   Difficulty      `json:"difficulty"`
	Category     Category        `json:"category"`
	CategoryTags []Category      `json:"category_tags"`
	Language     string          `json:"language"`
	DiffContent  string          `json:"diff_content"`
	FilePaths    []string        `json:"file_paths"`
	Metadata     json.RawMessage `json:"metadata"`
	IsPublished  bool            `json:"is_published"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// Validate checks that the exercise fields satisfy domain constraints.
func (e *Exercise) Validate() error {
	if e.Title == "" {
		return &ValidationError{Field: "title", Message: "must not be empty"}
	}
	if len(e.Title) > 200 {
		return &ValidationError{Field: "title", Message: "must be at most 200 characters"}
	}
	if e.Description == "" {
		return &ValidationError{Field: "description", Message: "must not be empty"}
	}
	if !e.Difficulty.Valid() {
		return &ValidationError{Field: "difficulty", Message: "invalid difficulty"}
	}
	if !e.Category.Valid() {
		return &ValidationError{Field: "category", Message: "invalid category"}
	}
	for i, tag := range e.CategoryTags {
		if !tag.Valid() {
			return &ValidationError{Field: "category_tags", Message: "invalid category at index " + itoa(i)}
		}
	}
	if e.Language == "" {
		return &ValidationError{Field: "language", Message: "must not be empty"}
	}
	if e.DiffContent == "" {
		return &ValidationError{Field: "diff_content", Message: "must not be empty"}
	}
	if len(e.DiffContent) > 100_000 {
		return &ValidationError{Field: "diff_content", Message: "must be at most 100,000 characters"}
	}
	return nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 4)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
