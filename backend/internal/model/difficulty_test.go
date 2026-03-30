package model_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestDifficultyValid(t *testing.T) {
	tests := []struct {
		difficulty model.Difficulty
		want       bool
	}{
		{model.DifficultyBeginner, true},
		{model.DifficultyIntermediate, true},
		{model.DifficultyAdvanced, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.difficulty), func(t *testing.T) {
			if got := tt.difficulty.Valid(); got != tt.want {
				t.Errorf("Difficulty(%q).Valid() = %v, want %v", tt.difficulty, got, tt.want)
			}
		})
	}
}

func TestAllDifficulties(t *testing.T) {
	diffs := model.AllDifficulties()
	if len(diffs) != 3 {
		t.Errorf("AllDifficulties() returned %d, want 3", len(diffs))
	}
}

func TestParseDifficultyValid(t *testing.T) {
	for _, raw := range []string{"beginner", "intermediate", "advanced"} {
		d, err := model.ParseDifficulty(raw)
		if err != nil {
			t.Errorf("ParseDifficulty(%q) returned unexpected error: %v", raw, err)
		}
		if d.String() != raw {
			t.Errorf("ParseDifficulty(%q).String() = %q", raw, d.String())
		}
	}
}

func TestParseDifficultyInvalid(t *testing.T) {
	_, err := model.ParseDifficulty("expert")
	if err == nil {
		t.Error("ParseDifficulty(\"expert\") expected error, got nil")
	}
}
