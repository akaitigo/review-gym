package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/akaitigo/review-gym/internal/model"
)

const (
	// exerciseListTTL is how long cached exercise lists live in Redis.
	exerciseListTTL = 5 * time.Minute
	// exerciseDetailTTL is how long a cached single exercise lives.
	exerciseDetailTTL = 10 * time.Minute
	// referenceReviewTTL is how long cached reference reviews live.
	referenceReviewTTL = 10 * time.Minute
)

// RedisCache wraps a primary store with Redis caching for read-heavy operations.
// It caches exercise listings, exercise details, and reference reviews.
// Write operations (comments, scores) are always delegated directly to the primary store.
type RedisCache struct {
	primary ExerciseStore
	reviews ReviewCommentStore
	refs    ReferenceReviewStore
	scores  ScoreStore
	rdb     *redis.Client
}

// NewRedisCache creates a RedisCache that wraps the given stores with Redis caching.
// It validates the Redis connection with a ping before returning.
func NewRedisCache(ctx context.Context, redisURL string, primary ExerciseStore, reviews ReviewCommentStore, refs ReferenceReviewStore, scores ScoreStore) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		closeErr := rdb.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("ping redis: %w (close error: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisCache{
		primary: primary,
		reviews: reviews,
		refs:    refs,
		scores:  scores,
		rdb:     rdb,
	}, nil
}

// Close closes the Redis client connection.
func (rc *RedisCache) Close() error {
	return rc.rdb.Close()
}

// exerciseListKey builds the cache key for an exercise list query.
func exerciseListKey(filter ExerciseFilter) string {
	key := "exercises:list"
	if filter.Category != nil {
		key += ":cat:" + string(*filter.Category)
	}
	if filter.Difficulty != nil {
		key += ":diff:" + string(*filter.Difficulty)
	}
	return key
}

// exerciseDetailKey builds the cache key for a single exercise.
func exerciseDetailKey(id string) string {
	return "exercises:detail:" + id
}

// referenceReviewKey builds the cache key for reference reviews.
func referenceReviewKey(exerciseID string) string {
	return "reference_reviews:" + exerciseID
}

// List returns exercises matching the given filter, using Redis cache.
func (rc *RedisCache) List(filter ExerciseFilter) ([]model.Exercise, error) {
	ctx := context.Background()
	cacheKey := exerciseListKey(filter)

	cached, err := rc.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var exercises []model.Exercise
		if jsonErr := json.Unmarshal(cached, &exercises); jsonErr == nil {
			return exercises, nil
		}
		// Cache corrupted, fall through to primary
	}

	exercises, err := rc.primary.List(filter)
	if err != nil {
		return nil, err
	}

	data, marshalErr := json.Marshal(exercises)
	if marshalErr == nil {
		rc.rdb.Set(ctx, cacheKey, data, exerciseListTTL)
	}

	return exercises, nil
}

// GetByID returns a single exercise by ID, using Redis cache.
func (rc *RedisCache) GetByID(id string) (*model.Exercise, error) {
	ctx := context.Background()
	cacheKey := exerciseDetailKey(id)

	cached, err := rc.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var exercise model.Exercise
		if jsonErr := json.Unmarshal(cached, &exercise); jsonErr == nil {
			return &exercise, nil
		}
	}

	exercise, err := rc.primary.GetByID(id)
	if err != nil {
		return nil, err
	}
	if exercise == nil {
		return nil, nil
	}

	data, marshalErr := json.Marshal(exercise)
	if marshalErr == nil {
		rc.rdb.Set(ctx, cacheKey, data, exerciseDetailTTL)
	}

	return exercise, nil
}

// Create delegates to the primary review comment store (no caching for writes).
func (rc *RedisCache) Create(comment *model.ReviewComment) error {
	return rc.reviews.Create(comment)
}

// ListByExerciseAndUser delegates to the primary review comment store.
func (rc *RedisCache) ListByExerciseAndUser(exerciseID, userID string) ([]model.ReviewComment, error) {
	return rc.reviews.ListByExerciseAndUser(exerciseID, userID)
}

// ListByExercise returns reference reviews, using Redis cache.
func (rc *RedisCache) ListByExercise(exerciseID string) ([]model.ReferenceReview, error) {
	ctx := context.Background()
	cacheKey := referenceReviewKey(exerciseID)

	cached, err := rc.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var reviews []model.ReferenceReview
		if jsonErr := json.Unmarshal(cached, &reviews); jsonErr == nil {
			return reviews, nil
		}
	}

	reviews, err := rc.refs.ListByExercise(exerciseID)
	if err != nil {
		return nil, err
	}

	data, marshalErr := json.Marshal(reviews)
	if marshalErr == nil {
		rc.rdb.Set(ctx, cacheKey, data, referenceReviewTTL)
	}

	return reviews, nil
}

// SaveScore delegates to the primary score store (no caching for writes).
func (rc *RedisCache) SaveScore(score *model.Score) error {
	return rc.scores.SaveScore(score)
}

// GetScoresByExerciseAndUser delegates to the primary score store.
func (rc *RedisCache) GetScoresByExerciseAndUser(exerciseID, userID string) ([]model.Score, error) {
	return rc.scores.GetScoresByExerciseAndUser(exerciseID, userID)
}

// GetScoresByUser delegates to the primary score store.
func (rc *RedisCache) GetScoresByUser(userID string) ([]model.Score, error) {
	return rc.scores.GetScoresByUser(userID)
}

// CountCompletedExercises delegates to the primary score store.
func (rc *RedisCache) CountCompletedExercises(userID string) (int, error) {
	return rc.scores.CountCompletedExercises(userID)
}
