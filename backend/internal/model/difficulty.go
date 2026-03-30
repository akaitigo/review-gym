package model

import "fmt"

// Difficulty represents the difficulty level of an exercise.
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

// AllDifficulties returns all valid difficulties in ascending order.
func AllDifficulties() []Difficulty {
	return []Difficulty{
		DifficultyBeginner,
		DifficultyIntermediate,
		DifficultyAdvanced,
	}
}

// Valid reports whether d is a recognized difficulty level.
func (d Difficulty) Valid() bool {
	switch d {
	case DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced:
		return true
	default:
		return false
	}
}

// String returns the string representation of the difficulty.
func (d Difficulty) String() string {
	return string(d)
}

// ParseDifficulty converts a raw string to a Difficulty, returning an error for unknown values.
func ParseDifficulty(raw string) (Difficulty, error) {
	d := Difficulty(raw)
	if !d.Valid() {
		return "", fmt.Errorf("unknown difficulty: %q", raw)
	}
	return d, nil
}
