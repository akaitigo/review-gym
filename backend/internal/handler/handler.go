// Package handler provides HTTP handlers for the review-gym API.
package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/akaitigo/review-gym/internal/store"
)

// defaultAnonymousUserID is used when X-User-ID header is not provided.
// This is a well-known UUID that represents an anonymous user.
const defaultAnonymousUserID = "00000000-0000-0000-0000-000000000000"

// uuidRegex validates UUID format (RFC 4122).
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

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

// resolveUserID extracts and validates the user ID from the X-User-ID header.
// If the header is empty, it returns the default anonymous UUID.
// If the header contains an invalid format, it returns an empty string and false.
func resolveUserID(r *http.Request) (string, bool) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return defaultAnonymousUserID, true
	}
	if !uuidRegex.MatchString(userID) {
		return "", false
	}
	return userID, true
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(resp)
}
