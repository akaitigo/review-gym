// Package store provides data access abstractions for the review-gym application.
package store

import (
	"context"

	"github.com/akaitigo/review-gym/internal/model"
)

// ExerciseFilter holds optional filter criteria for listing exercises.
type ExerciseFilter struct {
	Category   *model.Category
	Difficulty *model.Difficulty
}

// ExerciseStore defines operations on exercises.
type ExerciseStore interface {
	List(ctx context.Context, filter ExerciseFilter) ([]model.Exercise, error)
	GetByID(ctx context.Context, id string) (*model.Exercise, error)
}

// ReviewCommentStore defines operations on review comments.
type ReviewCommentStore interface {
	Create(ctx context.Context, comment *model.ReviewComment) error
	ListByExerciseAndUser(ctx context.Context, exerciseID, userID string) ([]model.ReviewComment, error)
}

// ReferenceReviewStore defines operations on reference reviews.
type ReferenceReviewStore interface {
	ListByExercise(ctx context.Context, exerciseID string) ([]model.ReferenceReview, error)
}

// ScoreStore defines operations on scoring results.
type ScoreStore interface {
	SaveScore(ctx context.Context, score *model.Score) error
	GetScoresByExerciseAndUser(ctx context.Context, exerciseID, userID string) ([]model.Score, error)
	GetScoresByUser(ctx context.Context, userID string) ([]model.Score, error)
	CountCompletedExercises(ctx context.Context, userID string) (int, error)
}
