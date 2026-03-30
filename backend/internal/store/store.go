// Package store provides data access abstractions for the review-gym application.
package store

import "github.com/akaitigo/review-gym/internal/model"

// ExerciseFilter holds optional filter criteria for listing exercises.
type ExerciseFilter struct {
	Category   *model.Category
	Difficulty *model.Difficulty
}

// ExerciseStore defines operations on exercises.
type ExerciseStore interface {
	List(filter ExerciseFilter) ([]model.Exercise, error)
	GetByID(id string) (*model.Exercise, error)
}

// ReviewCommentStore defines operations on review comments.
type ReviewCommentStore interface {
	Create(comment *model.ReviewComment) error
	ListByExerciseAndUser(exerciseID, userID string) ([]model.ReviewComment, error)
}

// ReferenceReviewStore defines operations on reference reviews.
type ReferenceReviewStore interface {
	ListByExercise(exerciseID string) ([]model.ReferenceReview, error)
}
