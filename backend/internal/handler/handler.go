// Package handler provides HTTP handlers for the review-gym API.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/review-gym/internal/store"
)

// Handler holds dependencies for all HTTP handlers.
type Handler struct {
	Exercises  store.ExerciseStore
	Reviews    store.ReviewCommentStore
	References store.ReferenceReviewStore
	Scores     store.ScoreStore
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(resp)
}
