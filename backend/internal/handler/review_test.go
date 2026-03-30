package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReview(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability detected",
		"category":    "security",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d; body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("expected id in response")
	}
	if resp["exercise_id"] != "00000001" {
		t.Errorf("expected exercise_id 00000001, got %v", resp["exercise_id"])
	}
}

func TestCreateReview_InvalidCategory(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := map[string]interface{}{
		"file_path":   "test.go",
		"line_number": 5,
		"content":     "Some comment",
		"category":    "invalid-category",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateReview_EmptyContent(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := map[string]interface{}{
		"file_path":   "test.go",
		"line_number": 5,
		"content":     "",
		"category":    "security",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateReview_ExerciseNotFound(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := map[string]interface{}{
		"file_path":   "test.go",
		"line_number": 5,
		"content":     "Some comment",
		"category":    "security",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/nonexistent/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestCreateReview_InvalidLineNumber(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 0,
		"content":     "Some comment",
		"category":    "security",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestListReviews(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// First create a review.
	body := map[string]interface{}{
		"file_path":   "internal/handler/user.go",
		"line_number": 21,
		"content":     "SQL injection vulnerability",
		"category":    "security",
	}
	b, _ := json.Marshal(body)

	createReq := httptest.NewRequest(http.MethodPost, "/api/exercises/00000001/reviews", bytes.NewReader(b))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "test-user")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("failed to create review: status %d", createRec.Code)
	}

	// Then list reviews.
	listReq := httptest.NewRequest(http.MethodGet, "/api/exercises/00000001/reviews", nil)
	listReq.Header.Set("X-User-ID", "test-user")
	listRec := httptest.NewRecorder()

	mux.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", listRec.Code)
	}

	var comments []map[string]interface{}
	if err := json.NewDecoder(listRec.Body).Decode(&comments); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
}

func TestListReviews_Empty(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/00000001/reviews", nil)
	req.Header.Set("X-User-ID", "test-user")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var comments []interface{}
	if err := json.NewDecoder(rec.Body).Decode(&comments); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestIsValidLineInDiff(t *testing.T) {
	diff := `--- a/test.go
+++ b/test.go
@@ -0,0 +1,5 @@
+package main
+
+func main() {
+	println("hello")
+}`

	tests := []struct {
		line int
		want bool
	}{
		{0, false},
		{1, true},
		{3, true},
		{5, true},
		{6, false},
		{-1, false},
	}

	for _, tt := range tests {
		got := isValidLineInDiff(diff, tt.line)
		if got != tt.want {
			t.Errorf("isValidLineInDiff(diff, %d) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestParseHunkNewStart(t *testing.T) {
	tests := []struct {
		header string
		want   int
	}{
		{"@@ -0,0 +1,35 @@", 1},
		{"@@ -10,5 +20,10 @@", 20},
		{"@@ -1 +1 @@", 1},
	}

	for _, tt := range tests {
		got := parseHunkNewStart(tt.header)
		if got != tt.want {
			t.Errorf("parseHunkNewStart(%q) = %d, want %d", tt.header, got, tt.want)
		}
	}
}
