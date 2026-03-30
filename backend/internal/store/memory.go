package store

import (
	"sync"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/seed"
)

// MemoryStore is an in-memory implementation of all store interfaces.
// Suitable for development and testing without a database.
type MemoryStore struct {
	mu               sync.RWMutex
	exercises        []model.Exercise
	referenceReviews map[string][]model.ReferenceReview // exerciseID -> reviews
	reviewComments   []model.ReviewComment
	nextCommentID    int
}

// NewMemoryStore creates a MemoryStore pre-populated with seed data.
func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		referenceReviews: make(map[string][]model.ReferenceReview),
		nextCommentID:    1,
	}
	ms.loadSeedData()
	return ms
}

func (ms *MemoryStore) loadSeedData() {
	all := seed.All()
	for i, ew := range all {
		id := generateID(i + 1)
		ex := ew.Exercise
		ex.ID = id
		now := time.Now()
		ex.CreatedAt = now
		ex.UpdatedAt = now
		ms.exercises = append(ms.exercises, ex)

		reviews := make([]model.ReferenceReview, len(ew.Reviews))
		for j, r := range ew.Reviews {
			r.ID = generateID((i+1)*100 + j + 1)
			r.ExerciseID = id
			r.CreatedAt = now
			reviews[j] = r
		}
		ms.referenceReviews[id] = reviews
	}
}

func generateID(n int) string {
	// Generate a zero-padded string ID for deterministic ordering.
	buf := make([]byte, 0, 8)
	if n == 0 {
		return "00000001"
	}
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	// Reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	// Pad to 8 characters
	for len(buf) < 8 {
		buf = append([]byte{'0'}, buf...)
	}
	return string(buf)
}

// List returns exercises matching the given filter.
func (ms *MemoryStore) List(filter ExerciseFilter) ([]model.Exercise, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.Exercise
	for _, ex := range ms.exercises {
		if !ex.IsPublished {
			continue
		}
		if filter.Category != nil && ex.Category != *filter.Category {
			// Also check category_tags
			found := false
			for _, tag := range ex.CategoryTags {
				if tag == *filter.Category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if filter.Difficulty != nil && ex.Difficulty != *filter.Difficulty {
			continue
		}
		result = append(result, ex)
	}
	return result, nil
}

// GetByID returns a single exercise by ID.
func (ms *MemoryStore) GetByID(id string) (*model.Exercise, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for i := range ms.exercises {
		if ms.exercises[i].ID == id {
			ex := ms.exercises[i]
			return &ex, nil
		}
	}
	return nil, nil
}

// Create stores a new review comment.
func (ms *MemoryStore) Create(comment *model.ReviewComment) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	comment.ID = generateID(ms.nextCommentID)
	ms.nextCommentID++
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	ms.reviewComments = append(ms.reviewComments, *comment)
	return nil
}

// ListByExerciseAndUser returns comments for a given exercise and user.
func (ms *MemoryStore) ListByExerciseAndUser(exerciseID, userID string) ([]model.ReviewComment, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.ReviewComment
	for _, c := range ms.reviewComments {
		if c.ExerciseID == exerciseID && c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

// ListByExercise returns all reference reviews for an exercise.
func (ms *MemoryStore) ListByExercise(exerciseID string) ([]model.ReferenceReview, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	reviews := ms.referenceReviews[exerciseID]
	result := make([]model.ReferenceReview, len(reviews))
	copy(result, reviews)
	return result, nil
}
