package handler

import "net/http"

// RegisterRoutes registers all API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/exercises", h.ListExercises)
	mux.HandleFunc("GET /api/exercises/{id}", h.GetExercise)
	mux.HandleFunc("POST /api/exercises/{id}/reviews", h.CreateReview)
	mux.HandleFunc("GET /api/exercises/{id}/reviews", h.ListReviews)
	mux.HandleFunc("POST /api/exercises/{id}/score", h.ScoreExercise)
	mux.HandleFunc("GET /api/exercises/{id}/scores", h.ListScores)
	mux.HandleFunc("GET /api/users/{id}/analytics", h.GetAnalytics)
	mux.HandleFunc("GET /api/users/{id}/recommendations", h.GetRecommendations)
}
