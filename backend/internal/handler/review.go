package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/review-gym/internal/model"
)

// createReviewRequest is the request body for creating a review comment.
type createReviewRequest struct {
	FilePath   string `json:"file_path"`
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
	Category   string `json:"category"`
}

// createReviewResponse is the response body after creating a review comment.
type createReviewResponse struct {
	ID         string         `json:"id"`
	ExerciseID string         `json:"exercise_id"`
	UserID     string         `json:"user_id"`
	FilePath   string         `json:"file_path"`
	LineNumber int            `json:"line_number"`
	Content    string         `json:"content"`
	Category   model.Category `json:"category"`
}

// CreateReview handles POST /api/exercises/{id}/reviews.
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

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

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cat, err := model.ParseCategory(req.Category)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category: must be one of security, performance, design, readability, error-handling")
		return
	}

	// Use a placeholder user ID until auth is implemented.
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	comment := &model.ReviewComment{
		ExerciseID: exerciseID,
		UserID:     userID,
		FilePath:   req.FilePath,
		LineNumber: req.LineNumber,
		Content:    req.Content,
		Category:   cat,
	}

	if err := comment.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate line number is within the diff range.
	if !isValidLineInDiff(exercise.DiffContent, req.LineNumber) {
		writeError(w, http.StatusBadRequest, "line_number is outside the valid diff range")
		return
	}

	if err := h.Reviews.Create(comment); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save review comment")
		return
	}

	resp := createReviewResponse{
		ID:         comment.ID,
		ExerciseID: comment.ExerciseID,
		UserID:     comment.UserID,
		FilePath:   comment.FilePath,
		LineNumber: comment.LineNumber,
		Content:    comment.Content,
		Category:   comment.Category,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// ListReviews handles GET /api/exercises/{id}/reviews.
func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.PathValue("id")
	if exerciseID == "" {
		writeError(w, http.StatusBadRequest, "missing exercise id")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	comments, err := h.Reviews.ListByExerciseAndUser(exerciseID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list reviews")
		return
	}
	if comments == nil {
		comments = []model.ReviewComment{}
	}

	writeJSON(w, http.StatusOK, comments)
}

// isValidLineInDiff checks if the given line number exists in the unified diff.
// It tracks lines in the new file (additions and context) using hunk headers
// to determine the valid line range for comments.
func isValidLineInDiff(diffContent string, lineNumber int) bool {
	if lineNumber <= 0 {
		return false
	}

	maxLine := 0
	currentLine := 0
	inHunk := false

	lines := splitLines(diffContent)
	for _, line := range lines {
		// Hunk header: @@ -old,count +new,count @@
		if len(line) >= 2 && line[0] == '@' && line[1] == '@' {
			inHunk = true
			currentLine = parseHunkNewStart(line)
			if currentLine > 0 {
				currentLine-- // Will be incremented on the first content line
			}
			continue
		}

		// Skip file headers (--- a/file, +++ b/file)
		if len(line) >= 3 && line[:3] == "---" {
			continue
		}
		if len(line) >= 3 && line[:3] == "+++" {
			continue
		}

		if !inHunk {
			continue
		}

		if len(line) > 0 && line[0] == '+' {
			currentLine++
			if currentLine > maxLine {
				maxLine = currentLine
			}
		} else if len(line) > 0 && line[0] == ' ' {
			currentLine++
			if currentLine > maxLine {
				maxLine = currentLine
			}
		}
		// Lines starting with '-' don't increment the new-file line number.
	}

	return lineNumber >= 1 && lineNumber <= maxLine
}

// splitLines splits text into lines, handling both \n and \r\n.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// parseHunkNewStart extracts the new-file start line from a unified diff hunk header.
// e.g., "@@ -0,0 +1,35 @@" returns 1.
func parseHunkNewStart(header string) int {
	// Find "+N" in the header.
	plusIdx := -1
	for i := 0; i < len(header)-1; i++ {
		if header[i] == '+' {
			plusIdx = i
			break
		}
	}
	if plusIdx < 0 {
		return 0
	}

	n := 0
	for i := plusIdx + 1; i < len(header); i++ {
		c := header[i]
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}
