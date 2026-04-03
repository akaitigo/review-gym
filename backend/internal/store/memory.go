package store

import (
	"sync"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/seed"
)

const (
	// maxReviewComments is the maximum number of review comments stored.
	// When exceeded, the oldest entries are evicted.
	maxReviewComments = 100_000
	// maxScores is the maximum number of score records stored.
	maxScores = 100_000
)

// MemoryStore is an in-memory implementation of all store interfaces.
// Suitable for development and testing without a database.
type MemoryStore struct {
	mu               sync.RWMutex
	exercises        []model.Exercise
	referenceReviews map[string][]model.ReferenceReview // exerciseID -> reviews
	reviewComments   []model.ReviewComment
	scores           []model.Score
	nextCommentID    int
	nextScoreID      int
}

// NewMemoryStore creates a MemoryStore pre-populated with seed data.
func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		referenceReviews: make(map[string][]model.ReferenceReview),
		nextCommentID:    1,
		nextScoreID:      1,
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
// If the store exceeds maxReviewComments, the oldest entries are evicted.
func (ms *MemoryStore) Create(comment *model.ReviewComment) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	comment.ID = generateID(ms.nextCommentID)
	ms.nextCommentID++
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	ms.reviewComments = append(ms.reviewComments, *comment)

	// Evict oldest entries if over capacity
	if len(ms.reviewComments) > maxReviewComments {
		excess := len(ms.reviewComments) - maxReviewComments
		ms.reviewComments = ms.reviewComments[excess:]
	}
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

// SaveScore persists a scoring result.
// It atomically computes the attempt number under the write lock to prevent
// race conditions where concurrent requests could get the same attempt number.
func (ms *MemoryStore) SaveScore(score *model.Score) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Recompute attempt number under lock to prevent concurrent duplicates.
	count := 0
	for _, s := range ms.scores {
		if s.ExerciseID == score.ExerciseID && s.UserID == score.UserID {
			count++
		}
	}
	score.AttemptNumber = count + 1

	score.ID = generateID(ms.nextScoreID)
	ms.nextScoreID++
	now := time.Now()
	score.CreatedAt = now
	ms.scores = append(ms.scores, *score)

	// Evict oldest entries if over capacity
	if len(ms.scores) > maxScores {
		excess := len(ms.scores) - maxScores
		ms.scores = ms.scores[excess:]
	}
	return nil
}

// GetScoresByExerciseAndUser returns all scores for a given exercise and user,
// ordered by attempt number.
func (ms *MemoryStore) GetScoresByExerciseAndUser(exerciseID, userID string) ([]model.Score, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.Score
	for _, s := range ms.scores {
		if s.ExerciseID == exerciseID && s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

// GetScoresByUser returns all scores for a given user.
func (ms *MemoryStore) GetScoresByUser(userID string) ([]model.Score, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.Score
	for _, s := range ms.scores {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

// CountCompletedExercises returns the number of distinct exercises scored by a user.
func (ms *MemoryStore) CountCompletedExercises(userID string) (int, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	seen := make(map[string]bool)
	for _, s := range ms.scores {
		if s.UserID == userID {
			seen[s.ExerciseID] = true
		}
	}
	return len(seen), nil
}
