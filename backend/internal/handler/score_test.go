package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScoreExercise(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// First, submit a review comment that matches a reference review.
	// Content must closely match the reference to achieve full precision.
	reviewBody := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability: user input is directly concatenated into the SQL query string.",
		"category":    "security",
	}
	b, _ := json.Marshal(reviewBody)

	createReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "test-user")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("failed to create review: status %d, body: %s", createRec.Code, createRec.Body.String())
	}

	// Now score the exercise.
	scoreReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
	scoreReq.Header.Set("X-User-ID", "test-user")
	scoreRec := httptest.NewRecorder()
	mux.ServeHTTP(scoreRec, scoreReq)

	if scoreRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", scoreRec.Code, scoreRec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(scoreRec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check required fields exist.
	for _, field := range []string{"id", "exercise_id", "user_id", "precision_score", "recall_score", "overall_score", "category_scores", "attempt_number", "matches", "missed_reviews", "false_positives"} {
		if _, ok := resp[field]; !ok {
			t.Errorf("missing field %q in response", field)
		}
	}

	// Precision should be 100 (1 comment, 1 match).
	precision, ok := resp["precision_score"].(float64)
	if !ok {
		t.Fatal("precision_score is not a number")
	}
	if precision != 100 {
		t.Errorf("precision = %.1f, want 100", precision)
	}

	// Recall should be > 0 (matched at least one reference).
	recall, ok := resp["recall_score"].(float64)
	if !ok {
		t.Fatal("recall_score is not a number")
	}
	if recall <= 0 {
		t.Errorf("recall = %.1f, want > 0", recall)
	}

	// Attempt number should be 1.
	attempt, ok := resp["attempt_number"].(float64)
	if !ok {
		t.Fatal("attempt_number is not a number")
	}
	if int(attempt) != 1 {
		t.Errorf("attempt = %d, want 1", int(attempt))
	}
}

func TestScoreExercise_NoComments(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
	req.Header.Set("X-User-ID", "test-user")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestScoreExercise_ExerciseNotFound(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/nonexistent/score", nil)
	req.Header.Set("X-User-ID", "test-user")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestScoreExercise_MultipleAttempts(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	reviewBody := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability: user input is directly concatenated into the SQL query string.",
		"category":    "security",
	}

	// Score twice, each time submitting a fresh comment before scoring.
	// After the fix, only comments from the current attempt are evaluated.
	for i := 1; i <= 2; i++ {
		b, _ := json.Marshal(reviewBody)
		createReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("X-User-ID", "test-user")
		createRec := httptest.NewRecorder()
		mux.ServeHTTP(createRec, createReq)

		if createRec.Code != http.StatusCreated {
			t.Fatalf("attempt %d: failed to create review: status %d", i, createRec.Code)
		}

		scoreReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
		scoreReq.Header.Set("X-User-ID", "test-user")
		scoreRec := httptest.NewRecorder()
		mux.ServeHTTP(scoreRec, scoreReq)

		if scoreRec.Code != http.StatusOK {
			t.Fatalf("attempt %d: expected status 200, got %d", i, scoreRec.Code)
		}

		var resp map[string]interface{}
		if err := json.NewDecoder(scoreRec.Body).Decode(&resp); err != nil {
			t.Fatalf("attempt %d: failed to decode: %v", i, err)
		}

		attempt := int(resp["attempt_number"].(float64))
		if attempt != i {
			t.Errorf("attempt %d: got attempt_number %d", i, attempt)
		}
	}
}

func TestScoreExercise_NoCommentsAfterPreviousScore(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Submit a comment and score it.
	reviewBody := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability",
		"category":    "security",
	}
	b, _ := json.Marshal(reviewBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "test-user")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	scoreReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
	scoreReq.Header.Set("X-User-ID", "test-user")
	scoreRec := httptest.NewRecorder()
	mux.ServeHTTP(scoreRec, scoreReq)

	if scoreRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on first score, got %d", scoreRec.Code)
	}

	// Score again without adding new comments — should get 400.
	scoreReq2 := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
	scoreReq2.Header.Set("X-User-ID", "test-user")
	scoreRec2 := httptest.NewRecorder()
	mux.ServeHTTP(scoreRec2, scoreReq2)

	if scoreRec2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when no new comments, got %d", scoreRec2.Code)
	}
}

func TestListScores(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Submit a comment and score.
	reviewBody := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability: user input is directly concatenated into the SQL query string.",
		"category":    "security",
	}
	b, _ := json.Marshal(reviewBody)

	createReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "test-user")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	scoreReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/score", nil)
	scoreReq.Header.Set("X-User-ID", "test-user")
	scoreRec := httptest.NewRecorder()
	mux.ServeHTTP(scoreRec, scoreReq)

	// List scores.
	listReq := httptest.NewRequest(http.MethodGet, "/api/exercises/00000001/scores", nil)
	listReq.Header.Set("X-User-ID", "test-user")
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", listRec.Code)
	}

	var scores []map[string]interface{}
	if err := json.NewDecoder(listRec.Body).Decode(&scores); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
}

func TestListScores_Empty(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/00000001/scores", nil)
	req.Header.Set("X-User-ID", "test-user")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var scores []interface{}
	if err := json.NewDecoder(rec.Body).Decode(&scores); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(scores) != 0 {
		t.Errorf("expected 0 scores, got %d", len(scores))
	}
}
