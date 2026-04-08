package handler

import (
	"net/http"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/scoring"
)

// scoreResponse is the response body for the scoring endpoint.
type scoreResponse struct {
	ID             string                `json:"id"`
	ExerciseID     string                `json:"exercise_id"`
	UserID         string                `json:"user_id"`
	PrecisionScore float64               `json:"precision_score"`
	RecallScore    float64               `json:"recall_score"`
	OverallScore   float64               `json:"overall_score"`
	CategoryScores []model.CategoryScore `json:"category_scores"`
	AttemptNumber  int                   `json:"attempt_number"`
	Matches        []scoring.Match       `json:"matches"`
	MissedReviews  []missedReviewItem    `json:"missed_reviews"`
	FalsePositives []falsePositiveItem   `json:"false_positives"`
}

// missedReviewItem is a simplified reference review for the response.
type missedReviewItem struct {
	FilePath    string         `json:"file_path"`
	LineNumber  int            `json:"line_number"`
	Content     string         `json:"content"`
	Category    model.Category `json:"category"`
	Severity    model.Severity `json:"severity"`
	Explanation string         `json:"explanation"`
}

// falsePositiveItem is a simplified user comment for the response.
type falsePositiveItem struct {
	FilePath   string         `json:"file_path"`
	LineNumber int            `json:"line_number"`
	Content    string         `json:"content"`
	Category   model.Category `json:"category"`
}

// ScoreExercise handles POST /api/exercises/{id}/score.
// It computes the score by comparing user comments against reference reviews.
func (h *Handler) ScoreExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.PathValue("id")
	if exerciseID == "" {
		writeError(w, http.StatusBadRequest, "missing exercise id")
		return
	}

	// Verify exercise exists.
	exercise, err := h.Exercises.GetByID(exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get exercise")
		return
	}
	if exercise == nil {
		writeError(w, http.StatusNotFound, "exercise not found")
		return
	}

	userID, ok := resolveUserID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "X-User-ID must be a valid UUID")
		return
	}

	// Fetch existing scores to filter comments by attempt boundary.
	existingScores, err := h.Scores.GetScoresByExerciseAndUser(exerciseID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get existing scores")
		return
	}

	// Get user's review comments for this exercise.
	allComments, err := h.Reviews.ListByExerciseAndUser(exerciseID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list review comments")
		return
	}

	// Filter to only comments from the current attempt (i.e., after the last score).
	// This prevents previously scored comments from inflating the current attempt.
	var comments []model.ReviewComment
	if len(existingScores) > 0 {
		lastScoreTime := existingScores[len(existingScores)-1].CreatedAt
		for _, c := range allComments {
			if c.CreatedAt.After(lastScoreTime) {
				comments = append(comments, c)
			}
		}
	} else {
		comments = allComments
	}

	if len(comments) == 0 {
		writeError(w, http.StatusBadRequest, "no review comments found; submit at least one comment before scoring")
		return
	}

	// Get reference reviews for this exercise.
	references, err := h.References.ListByExercise(exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list reference reviews")
		return
	}

	// Compute the score.
	result := scoring.Compute(comments, references)

	// Persist the score. Pass 0 as attemptNumber — the store layer
	// atomically computes the correct attempt number.
	score, err := result.ToScore(userID, exerciseID, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create score record")
		return
	}

	if err := h.Scores.SaveScore(score); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save score")
		return
	}

	// Build response.
	missed := make([]missedReviewItem, 0, len(result.MissedReviews))
	for _, m := range result.MissedReviews {
		missed = append(missed, missedReviewItem{
			FilePath:    m.FilePath,
			LineNumber:  m.LineNumber,
			Content:     m.Content,
			Category:    m.Category,
			Severity:    m.Severity,
			Explanation: m.Explanation,
		})
	}

	fps := make([]falsePositiveItem, 0, len(result.FalsePositives))
	for _, fp := range result.FalsePositives {
		fps = append(fps, falsePositiveItem{
			FilePath:   fp.FilePath,
			LineNumber: fp.LineNumber,
			Content:    fp.Content,
			Category:   fp.Category,
		})
	}

	resp := scoreResponse{
		ID:             score.ID,
		ExerciseID:     exerciseID,
		UserID:         userID,
		PrecisionScore: result.PrecisionScore,
		RecallScore:    result.RecallScore,
		OverallScore:   result.OverallScore,
		CategoryScores: result.CategoryScores,
		AttemptNumber:  score.AttemptNumber,
		Matches:        result.Matches,
		MissedReviews:  missed,
		FalsePositives: fps,
	}

	// Ensure non-nil slices for JSON.
	if resp.Matches == nil {
		resp.Matches = []scoring.Match{}
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListScores handles GET /api/exercises/{id}/scores.
// Returns the score history for the current user on a specific exercise.
func (h *Handler) ListScores(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.PathValue("id")
	if exerciseID == "" {
		writeError(w, http.StatusBadRequest, "missing exercise id")
		return
	}

	userID, ok := resolveUserID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "X-User-ID must be a valid UUID")
		return
	}

	scores, err := h.Scores.GetScoresByExerciseAndUser(exerciseID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list scores")
		return
	}
	if scores == nil {
		scores = []model.Score{}
	}

	writeJSON(w, http.StatusOK, scores)
}
