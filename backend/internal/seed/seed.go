// Package seed provides exercise and reference review seed data for the review-gym database.
package seed

import (
	"github.com/akaitigo/review-gym/internal/model"
)

// ExerciseWithReviews pairs an exercise with its reference reviews.
type ExerciseWithReviews struct {
	Exercise model.Exercise
	Reviews  []model.ReferenceReview
}

// All returns all seed exercises with their reference reviews.
// Each exercise represents an anonymized OSS pull request with curated model reviews.
func All() []ExerciseWithReviews {
	return []ExerciseWithReviews{
		exercise01SQLInjection(),
		exercise02UnboundedQuery(),
		exercise03GodFunction(),
		exercise04ErrorSwallowing(),
		exercise05HardcodedSecret(),
		exercise06NPlus1Query(),
		exercise07RaceCondition(),
		exercise08XSSVulnerability(),
		exercise09MemoryLeak(),
		exercise10MagicNumbers(),
		exercise11InsecureDeserialization(),
		exercise12UnbufferedChannel(),
	}
}
