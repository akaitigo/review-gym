package store

import (
	"context"
	"fmt"
	"io"
	"log"
)

// StoreType represents which backing store implementation to use.
type StoreType string

const (
	// StoreTypeMemory uses in-memory storage (default, data lost on restart).
	StoreTypeMemory StoreType = "memory"
	// StoreTypePostgres uses PostgreSQL for persistent storage.
	StoreTypePostgres StoreType = "postgres"
)

// Stores holds all store interface implementations and any closeable resources.
type Stores struct {
	Exercises  ExerciseStore
	Reviews    ReviewCommentStore
	References ReferenceReviewStore
	Scores     ScoreStore
	closers    []io.Closer
}

// Close releases all underlying resources (database connections, Redis, etc.).
func (s *Stores) Close() {
	for _, c := range s.closers {
		if err := c.Close(); err != nil {
			log.Printf("store close error: %v", err)
		}
	}
}

// Config holds configuration for creating stores.
type Config struct {
	// StoreType selects the backing store: "memory" or "postgres".
	StoreType StoreType
	// DatabaseURL is the PostgreSQL connection string (required when StoreType is "postgres").
	DatabaseURL string
	// RedisURL is the Redis connection string (optional, enables caching layer).
	RedisURL string
}

// NewStores creates and initializes all stores based on the given configuration.
// When StoreType is "postgres", it connects to PostgreSQL.
// When RedisURL is set, it wraps read-heavy stores with a Redis cache layer.
// Returns Stores that must be closed when no longer needed.
func NewStores(ctx context.Context, cfg Config) (*Stores, error) {
	stores := &Stores{}

	switch cfg.StoreType {
	case StoreTypePostgres:
		if cfg.DatabaseURL == "" {
			return nil, fmt.Errorf("DATABASE_URL is required when STORE_TYPE is postgres")
		}

		ps, err := NewPostgresStore(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("create postgres store: %w", err)
		}
		stores.closers = append(stores.closers, ps)

		stores.Exercises = ps
		stores.Reviews = ps
		stores.References = ps
		stores.Scores = ps

		log.Println("store: using PostgreSQL")

	case StoreTypeMemory, "":
		ms := NewMemoryStore()
		stores.Exercises = ms
		stores.Reviews = ms
		stores.References = ms
		stores.Scores = ms

		log.Println("store: using in-memory (data will be lost on restart)")

	default:
		return nil, fmt.Errorf("unknown STORE_TYPE: %q (valid: memory, postgres)", cfg.StoreType)
	}

	// Wrap with Redis cache if configured.
	if cfg.RedisURL != "" {
		rc, err := NewRedisCache(ctx, cfg.RedisURL, stores.Exercises, stores.Reviews, stores.References, stores.Scores)
		if err != nil {
			log.Printf("store: Redis cache unavailable, continuing without cache: %v", err)
		} else {
			stores.closers = append(stores.closers, rc)
			stores.Exercises = rc
			stores.References = rc
			// Reviews and Scores are delegated through RedisCache already
			stores.Reviews = rc
			stores.Scores = rc

			log.Println("store: Redis cache enabled")
		}
	} else {
		log.Println("store: Redis cache disabled (REDIS_URL not set)")
	}

	return stores, nil
}
