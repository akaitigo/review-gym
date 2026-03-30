package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/store"
)

func setupAnalyticsHandler() (*Handler, *store.MemoryStore) {
	ms := store.NewMemoryStore()
	h := &Handler{
		Exercises:  ms,
		Reviews:    ms,
		References: ms,
		Scores:     ms,
	}
	return h, ms
}

func seedScores(t *testing.T, ms *store.MemoryStore, userID string, count int) {
	t.Helper()

	exercises, err := ms.List(store.ExerciseFilter{})
	if err != nil {
		t.Fatalf("failed to list exercises: %v", err)
	}

	for i := 0; i < count && i < len(exercises); i++ {
		catScores := []model.CategoryScore{
			{Category: model.CategorySecurity, Score: float64(40 + i*10), MaxPoints: 3, Earned: float64(1 + i)},
			{Category: model.CategoryPerformance, Score: float64(60 + i*5), MaxPoints: 2, Earned: float64(1 + i)},
			{Category: model.CategoryDesign, Score: float64(80 - i*5), MaxPoints: 2, Earned: float64(1)},
			{Category: model.CategoryReadability, Score: float64(70 + i*3), MaxPoints: 1, Earned: float64(1)},
			{Category: model.CategoryErrorHandling, Score: float64(30 + i*2), MaxPoints: 3, Earned: float64(1)},
		}
		catJSON, err := json.Marshal(catScores)
		if err != nil {
			t.Fatalf("failed to marshal category scores: %v", err)
		}

		score := &model.Score{
			UserID:         userID,
			ExerciseID:     exercises[i].ID,
			PrecisionScore: float64(50 + i*10),
			RecallScore:    float64(60 + i*5),
			OverallScore:   float64(55 + i*7),
			CategoryScores: catJSON,
			AttemptNumber:  1,
		}
		if err := ms.SaveScore(score); err != nil {
			t.Fatalf("failed to save score: %v", err)
		}
	}
}

func TestGetAnalytics_MissingUserID(t *testing.T) {
	h, _ := setupAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users//analytics", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	h.GetAnalytics(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetAnalytics_NoScores(t *testing.T) {
	h, _ := setupAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users/unknown-user/analytics", nil)
	req.SetPathValue("id", "unknown-user")
	w := httptest.NewRecorder()

	h.GetAnalytics(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetAnalytics_InsufficientExercises(t *testing.T) {
	h, ms := setupAnalyticsHandler()
	seedScores(t, ms, "user-1", 2)

	req := httptest.NewRequest(http.MethodGet, "/api/users/user-1/analytics", nil)
	req.SetPathValue("id", "user-1")
	w := httptest.NewRecorder()

	h.GetAnalytics(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	msg, ok := resp["message"].(string)
	if !ok || msg == "" {
		t.Error("expected a message about insufficient exercises")
	}
}

func TestGetAnalytics_Success(t *testing.T) {
	h, ms := setupAnalyticsHandler()
	seedScores(t, ms, "user-1", 4)

	req := httptest.NewRequest(http.MethodGet, "/api/users/user-1/analytics", nil)
	req.SetPathValue("id", "user-1")
	w := httptest.NewRecorder()

	h.GetAnalytics(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp analyticsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.UserID != "user-1" {
		t.Errorf("expected user_id 'user-1', got %q", resp.UserID)
	}

	if resp.TotalExercisesCompleted != 4 {
		t.Errorf("expected 4 completed exercises, got %d", resp.TotalExercisesCompleted)
	}

	if resp.TotalAttempts != 4 {
		t.Errorf("expected 4 attempts, got %d", resp.TotalAttempts)
	}

	if len(resp.Categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(resp.Categories))
	}

	if resp.OverallAverageScore <= 0 {
		t.Error("expected positive overall average score")
	}

	if len(resp.ScoreHistory) != 4 {
		t.Errorf("expected 4 score history points, got %d", len(resp.ScoreHistory))
	}

	// At least error-handling should be a weakness (scores: 30, 32, 34, 36 => avg 33).
	hasWeakness := false
	for _, cat := range resp.Categories {
		if cat.Category == model.CategoryErrorHandling && cat.IsWeakness {
			hasWeakness = true
		}
	}
	if !hasWeakness {
		t.Error("expected error-handling to be flagged as weakness")
	}
}

func TestGetRecommendations_MissingUserID(t *testing.T) {
	h, _ := setupAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users//recommendations", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	h.GetRecommendations(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetRecommendations_NoScores(t *testing.T) {
	h, _ := setupAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users/unknown/recommendations", nil)
	req.SetPathValue("id", "unknown")
	w := httptest.NewRecorder()

	h.GetRecommendations(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetRecommendations_InsufficientExercises(t *testing.T) {
	h, ms := setupAnalyticsHandler()
	seedScores(t, ms, "user-2", 2)

	req := httptest.NewRequest(http.MethodGet, "/api/users/user-2/recommendations", nil)
	req.SetPathValue("id", "user-2")
	w := httptest.NewRecorder()

	h.GetRecommendations(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	msg, ok := resp["message"].(string)
	if !ok || msg == "" {
		t.Error("expected a message about insufficient exercises")
	}
}

func TestGetRecommendations_Success(t *testing.T) {
	h, ms := setupAnalyticsHandler()
	seedScores(t, ms, "user-3", 4)

	req := httptest.NewRequest(http.MethodGet, "/api/users/user-3/recommendations", nil)
	req.SetPathValue("id", "user-3")
	w := httptest.NewRecorder()

	h.GetRecommendations(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp recommendationsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.UserID != "user-3" {
		t.Errorf("expected user_id 'user-3', got %q", resp.UserID)
	}

	// Should have weakness categories based on the seeded scores.
	if len(resp.WeaknessCategories) == 0 {
		t.Error("expected at least one weakness category")
	}

	// Recommendations should be ordered: unattempted first.
	for i := 0; i < len(resp.Recommendations)-1; i++ {
		if resp.Recommendations[i].PreviouslyAttempted && !resp.Recommendations[i+1].PreviouslyAttempted {
			t.Error("expected unattempted exercises to come before attempted ones")
		}
	}
}

func TestComputeTrend(t *testing.T) {
	tests := []struct {
		name     string
		scores   []float64
		expected string
	}{
		{
			name:     "single score",
			scores:   []float64{50},
			expected: "stagnating",
		},
		{
			name:     "improving",
			scores:   []float64{30, 40, 70, 80},
			expected: "improving",
		},
		{
			name:     "declining",
			scores:   []float64{80, 70, 30, 20},
			expected: "declining",
		},
		{
			name:     "stagnating",
			scores:   []float64{50, 52, 50, 51},
			expected: "stagnating",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeTrend(tt.scores)
			if result != tt.expected {
				t.Errorf("computeTrend(%v) = %q, want %q", tt.scores, result, tt.expected)
			}
		})
	}
}

func TestComputeConsecutiveDays(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	tests := []struct {
		name     string
		dates    []string
		expected int
	}{
		{
			name:     "empty",
			dates:    []string{},
			expected: 0,
		},
		{
			name:     "today only",
			dates:    []string{today},
			expected: 1,
		},
		{
			name:     "today and yesterday",
			dates:    []string{yesterday, today},
			expected: 2,
		},
		{
			name:     "three day streak ending today",
			dates:    []string{twoDaysAgo, yesterday, today},
			expected: 3,
		},
		{
			name:     "streak ending yesterday",
			dates:    []string{threeDaysAgo, twoDaysAgo, yesterday},
			expected: 3,
		},
		{
			name:     "old dates only",
			dates:    []string{"2020-01-01", "2020-01-02"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeConsecutiveDays(tt.dates)
			if result != tt.expected {
				t.Errorf("computeConsecutiveDays(%v) = %d, want %d", tt.dates, result, tt.expected)
			}
		})
	}
}
