package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/review-gym/internal/store"
)

func newTestHandler() *Handler {
	ms := store.NewMemoryStore()
	return &Handler{
		Exercises:  ms,
		Reviews:    ms,
		References: ms,
	}
}

func TestListExercises(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var items []json.RawMessage
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected at least one exercise in response")
	}

	// Verify the response does not contain diff_content (list should omit it).
	var firstItem map[string]interface{}
	if err := json.Unmarshal(items[0], &firstItem); err != nil {
		t.Fatalf("failed to unmarshal first item: %v", err)
	}
	if _, hasDiff := firstItem["diff_content"]; hasDiff {
		t.Error("list response should not contain diff_content")
	}
}

func TestListExercises_FilterByCategory(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises?category=security", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var items []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected at least one security exercise")
	}
}

func TestListExercises_InvalidCategory(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises?category=invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestListExercises_FilterByDifficulty(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises?difficulty=beginner", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var items []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for _, item := range items {
		if item["difficulty"] != "beginner" {
			t.Errorf("expected difficulty beginner, got %v", item["difficulty"])
		}
	}
}

func TestListExercises_InvalidDifficulty(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises?difficulty=invalid", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestGetExercise(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// First, get a valid ID from the list.
	listReq := httptest.NewRequest(http.MethodGet, "/api/exercises", nil)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)

	var items []map[string]interface{}
	if err := json.NewDecoder(listRec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("no exercises found")
	}

	id, ok := items[0]["id"].(string)
	if !ok {
		t.Fatal("id is not a string")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/"+id, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var exercise map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&exercise); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if exercise["diff_content"] == nil || exercise["diff_content"] == "" {
		t.Error("expected diff_content in detail response")
	}
}

func TestGetExercise_NotFound(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/nonexistent", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
