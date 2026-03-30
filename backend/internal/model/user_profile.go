package model

import "time"

// UserProfile represents a registered user of the review-gym platform.
type UserProfile struct {
	ID                      string     `json:"id"`
	DisplayName             string     `json:"display_name"`
	Email                   string     `json:"email,omitempty"`
	GitHubID                string     `json:"github_id,omitempty"`
	AvatarURL               string     `json:"avatar_url,omitempty"`
	WeaknessCategories      []Category `json:"weakness_categories"`
	TotalExercisesCompleted int        `json:"total_exercises_completed"`
	ConsecutiveDays         int        `json:"consecutive_days"`
	LastPracticeAt          *time.Time `json:"last_practice_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// Validate checks that the user profile fields satisfy domain constraints.
func (u *UserProfile) Validate() error {
	if u.DisplayName == "" {
		return &ValidationError{Field: "display_name", Message: "must not be empty"}
	}
	if len(u.DisplayName) > 100 {
		return &ValidationError{Field: "display_name", Message: "must be at most 100 characters"}
	}
	for i, cat := range u.WeaknessCategories {
		if !cat.Valid() {
			return &ValidationError{Field: "weakness_categories", Message: "invalid category at index " + itoa(i)}
		}
	}
	return nil
}
