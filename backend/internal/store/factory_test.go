package store

import (
	"context"
	"testing"
	"time"
)

func TestNewStores_MemoryDefault(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stores, err := NewStores(ctx, Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stores.Close()

	if stores.Exercises == nil {
		t.Error("Exercises store is nil")
	}
	if stores.Reviews == nil {
		t.Error("Reviews store is nil")
	}
	if stores.References == nil {
		t.Error("References store is nil")
	}
	if stores.Scores == nil {
		t.Error("Scores store is nil")
	}

	// Verify it loaded seed data.
	exercises, err := stores.Exercises.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Error("expected seed exercises to be loaded")
	}
}

func TestNewStores_MemoryExplicit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stores, err := NewStores(ctx, Config{StoreType: StoreTypeMemory})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stores.Close()

	exercises, err := stores.Exercises.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Error("expected seed exercises to be loaded")
	}
}

func TestNewStores_PostgresRequiresDatabaseURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewStores(ctx, Config{StoreType: StoreTypePostgres})
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func TestNewStores_InvalidStoreType(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewStores(ctx, Config{StoreType: "invalid"})
	if err == nil {
		t.Fatal("expected error for invalid store type")
	}
}

func TestNewStores_RedisFallbackOnInvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// An invalid Redis URL should log a warning and continue without cache.
	stores, err := NewStores(ctx, Config{
		StoreType: StoreTypeMemory,
		RedisURL:  "redis://invalid-host-that-does-not-exist:12345",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stores.Close()

	// Should still work with memory store.
	exercises, err := stores.Exercises.List(ctx, ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exercises) == 0 {
		t.Error("expected seed exercises to be loaded")
	}
}
