package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/store"
)

// weaknessThresholdPct is the percentage below the overall average
// at which a category is considered a weakness (80%).
const weaknessThresholdPct = 0.80

// minExercisesForAnalytics is the minimum number of completed exercises
// required before analytics data is available.
const minExercisesForAnalytics = 3

// categoryAnalytics represents the analytics for a single review category.
type categoryAnalytics struct {
	Category     model.Category `json:"category"`
	AverageScore float64        `json:"average_score"`
	MinScore     float64        `json:"min_score"`
	MaxScore     float64        `json:"max_score"`
	Trend        string         `json:"trend"` // "improving", "stagnating", "declining"
	IsWeakness   bool           `json:"is_weakness"`
}

// analyticsResponse is the response body for the analytics endpoint.
type analyticsResponse struct {
	UserID                  string              `json:"user_id"`
	TotalExercisesCompleted int                 `json:"total_exercises_completed"`
	TotalAttempts           int                 `json:"total_attempts"`
	OverallAverageScore     float64             `json:"overall_average_score"`
	Categories              []categoryAnalytics `json:"categories"`
	WeaknessCategories      []model.Category    `json:"weakness_categories"`
	ScoreHistory            []scoreHistoryPoint `json:"score_history"`
	ConsecutiveDays         int                 `json:"consecutive_days"`
}

// scoreHistoryPoint represents a single data point for the score trend chart.
type scoreHistoryPoint struct {
	Date         string  `json:"date"`
	OverallScore float64 `json:"overall_score"`
	AttemptIndex int     `json:"attempt_index"`
}

// GetAnalytics handles GET /api/users/{id}/analytics.
// It computes category-level averages, weakness detection, and score trends.
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	scores, err := h.Scores.GetScoresByUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get scores")
		return
	}

	if len(scores) == 0 {
		writeError(w, http.StatusNotFound, "user not found or no scores recorded")
		return
	}

	completedCount, err := h.Scores.CountCompletedExercises(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count exercises")
		return
	}

	if completedCount < minExercisesForAnalytics {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"user_id":                   userID,
			"total_exercises_completed": completedCount,
			"total_attempts":            len(scores),
			"message":                   "complete at least 3 exercises to unlock analytics",
			"min_exercises_required":    minExercisesForAnalytics,
		})
		return
	}

	// Parse category scores from each score record.
	type catAccumulator struct {
		scores []float64
	}
	catAcc := make(map[model.Category]*catAccumulator)
	for _, cat := range model.AllCategories() {
		catAcc[cat] = &catAccumulator{}
	}

	var overallSum float64
	for _, s := range scores {
		overallSum += s.OverallScore

		var catScores []model.CategoryScore
		if err := json.Unmarshal(s.CategoryScores, &catScores); err != nil {
			continue
		}
		for _, cs := range catScores {
			if acc, ok := catAcc[cs.Category]; ok {
				acc.scores = append(acc.scores, cs.Score)
			}
		}
	}

	overallAvg := overallSum / float64(len(scores))
	weaknessThreshold := overallAvg * weaknessThresholdPct

	var categories []categoryAnalytics
	var weaknesses []model.Category

	for _, cat := range model.AllCategories() {
		acc := catAcc[cat]
		ca := categoryAnalytics{
			Category: cat,
		}

		if len(acc.scores) > 0 {
			var sum, min, max float64
			min = acc.scores[0]
			max = acc.scores[0]
			for _, v := range acc.scores {
				sum += v
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}
			ca.AverageScore = roundTo1(sum / float64(len(acc.scores)))
			ca.MinScore = roundTo1(min)
			ca.MaxScore = roundTo1(max)
			ca.Trend = computeTrend(acc.scores)
			ca.IsWeakness = ca.AverageScore < roundTo1(weaknessThreshold)
		}

		if ca.IsWeakness {
			weaknesses = append(weaknesses, cat)
		}
		categories = append(categories, ca)
	}

	// Build score history (chronological by CreatedAt).
	sortedScores := make([]model.Score, len(scores))
	copy(sortedScores, scores)
	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].CreatedAt.Before(sortedScores[j].CreatedAt)
	})

	history := make([]scoreHistoryPoint, 0, len(sortedScores))
	for i, s := range sortedScores {
		history = append(history, scoreHistoryPoint{
			Date:         s.CreatedAt.Format(time.DateOnly),
			OverallScore: s.OverallScore,
			AttemptIndex: i + 1,
		})
	}

	// Compute consecutive practice days.
	scoreDates, err := h.Scores.GetScoreDates(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get score dates")
		return
	}
	consecutive := computeConsecutiveDays(scoreDates)

	resp := analyticsResponse{
		UserID:                  userID,
		TotalExercisesCompleted: completedCount,
		TotalAttempts:           len(scores),
		OverallAverageScore:     roundTo1(overallAvg),
		Categories:              categories,
		WeaknessCategories:      weaknesses,
		ScoreHistory:            history,
		ConsecutiveDays:         consecutive,
	}

	if resp.WeaknessCategories == nil {
		resp.WeaknessCategories = []model.Category{}
	}

	writeJSON(w, http.StatusOK, resp)
}

// computeTrend determines the trend direction for a series of scores.
// It compares the average of the first half with the average of the second half.
func computeTrend(scores []float64) string {
	if len(scores) < 2 {
		return "stagnating"
	}

	mid := len(scores) / 2
	firstHalf := scores[:mid]
	secondHalf := scores[mid:]

	var firstAvg, secondAvg float64
	for _, v := range firstHalf {
		firstAvg += v
	}
	firstAvg /= float64(len(firstHalf))

	for _, v := range secondHalf {
		secondAvg += v
	}
	secondAvg /= float64(len(secondHalf))

	diff := secondAvg - firstAvg
	if diff > 5 {
		return "improving"
	}
	if diff < -5 {
		return "declining"
	}
	return "stagnating"
}

// computeConsecutiveDays calculates the current streak of consecutive practice days.
func computeConsecutiveDays(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	// Sort dates in descending order for streak calculation.
	sorted := make([]string, len(dates))
	copy(sorted, dates)
	sort.Sort(sort.Reverse(sort.StringSlice(sorted)))

	today := time.Now().Format("2006-01-02")

	// The streak must include today or yesterday.
	streak := 0
	expected := today

	for _, d := range sorted {
		if d == expected {
			streak++
			t, err := time.Parse("2006-01-02", expected)
			if err != nil {
				break
			}
			expected = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else if streak == 0 {
			// Allow yesterday as the start of streak.
			yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			if d == yesterday {
				streak++
				t, err := time.Parse("2006-01-02", yesterday)
				if err != nil {
					break
				}
				expected = t.AddDate(0, 0, -1).Format("2006-01-02")
			} else {
				break
			}
		} else {
			break
		}
	}
	return streak
}

// roundTo1 rounds a float64 to 1 decimal place.
func roundTo1(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}

// recommendationItem represents a recommended exercise for the user.
type recommendationItem struct {
	Exercise            model.Exercise `json:"exercise"`
	RecommendedReason   string         `json:"recommended_reason"`
	TargetWeakness      model.Category `json:"target_weakness"`
	PreviouslyAttempted bool           `json:"previously_attempted"`
}

// recommendationsResponse is the response body for the recommendations endpoint.
type recommendationsResponse struct {
	UserID             string               `json:"user_id"`
	WeaknessCategories []model.Category     `json:"weakness_categories"`
	Recommendations    []recommendationItem `json:"recommendations"`
}

// GetRecommendations handles GET /api/users/{id}/recommendations.
// It returns exercises recommended based on the user's weakness categories.
func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	scores, err := h.Scores.GetScoresByUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get scores")
		return
	}

	if len(scores) == 0 {
		writeError(w, http.StatusNotFound, "user not found or no scores recorded")
		return
	}

	completedCount, err := h.Scores.CountCompletedExercises(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count exercises")
		return
	}

	if completedCount < minExercisesForAnalytics {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"user_id":                   userID,
			"total_exercises_completed": completedCount,
			"message":                   "complete at least 3 exercises to unlock recommendations",
			"min_exercises_required":    minExercisesForAnalytics,
			"recommendations":           []recommendationItem{},
		})
		return
	}

	// Determine weakness categories.
	catAvg := computeCategoryAverages(scores)
	var overallSum float64
	for _, s := range scores {
		overallSum += s.OverallScore
	}
	overallAvg := overallSum / float64(len(scores))
	weaknessThreshold := overallAvg * weaknessThresholdPct

	var weaknesses []model.Category
	for _, cat := range model.AllCategories() {
		if avg, ok := catAvg[cat]; ok && avg < weaknessThreshold {
			weaknesses = append(weaknesses, cat)
		}
	}

	// Get completed exercise IDs.
	completedIDs, err := h.Scores.GetCompletedExerciseIDs(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get completed exercises")
		return
	}

	// Get all exercises.
	allExercises, err := h.Exercises.List(store.ExerciseFilter{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list exercises")
		return
	}

	// Build recommendations: exercises matching weakness categories,
	// prioritizing unattempted exercises and ascending difficulty.
	var recs []recommendationItem
	for _, weakness := range weaknesses {
		for _, ex := range allExercises {
			if !exerciseMatchesCategory(ex, weakness) {
				continue
			}
			attempted := completedIDs[ex.ID]
			reason := "Targets your weakness in " + string(weakness)
			if !attempted {
				reason += " (not yet attempted)"
			}
			recs = append(recs, recommendationItem{
				Exercise:            ex,
				RecommendedReason:   reason,
				TargetWeakness:      weakness,
				PreviouslyAttempted: attempted,
			})
		}
	}

	// Sort: unattempted first, then by difficulty (beginner -> advanced).
	sort.Slice(recs, func(i, j int) bool {
		// Unattempted exercises first.
		if recs[i].PreviouslyAttempted != recs[j].PreviouslyAttempted {
			return !recs[i].PreviouslyAttempted
		}
		// Then by difficulty ascending.
		return difficultyOrder(recs[i].Exercise.Difficulty) < difficultyOrder(recs[j].Exercise.Difficulty)
	})

	// Deduplicate by exercise ID (keep first occurrence).
	seen := make(map[string]bool)
	var deduped []recommendationItem
	for _, rec := range recs {
		if !seen[rec.Exercise.ID] {
			seen[rec.Exercise.ID] = true
			deduped = append(deduped, rec)
		}
	}

	// Limit to 10 recommendations.
	if len(deduped) > 10 {
		deduped = deduped[:10]
	}

	if weaknesses == nil {
		weaknesses = []model.Category{}
	}
	if deduped == nil {
		deduped = []recommendationItem{}
	}

	resp := recommendationsResponse{
		UserID:             userID,
		WeaknessCategories: weaknesses,
		Recommendations:    deduped,
	}

	writeJSON(w, http.StatusOK, resp)
}

// computeCategoryAverages extracts category score averages from a user's score records.
func computeCategoryAverages(scores []model.Score) map[model.Category]float64 {
	type acc struct {
		sum   float64
		count int
	}
	catAcc := make(map[model.Category]*acc)

	for _, s := range scores {
		var catScores []model.CategoryScore
		if err := json.Unmarshal(s.CategoryScores, &catScores); err != nil {
			continue
		}
		for _, cs := range catScores {
			a, ok := catAcc[cs.Category]
			if !ok {
				a = &acc{}
				catAcc[cs.Category] = a
			}
			a.sum += cs.Score
			a.count++
		}
	}

	result := make(map[model.Category]float64)
	for cat, a := range catAcc {
		if a.count > 0 {
			result[cat] = a.sum / float64(a.count)
		}
	}
	return result
}

// exerciseMatchesCategory checks if an exercise matches a given category
// (either primary category or category tags).
func exerciseMatchesCategory(ex model.Exercise, cat model.Category) bool {
	if ex.Category == cat {
		return true
	}
	for _, tag := range ex.CategoryTags {
		if tag == cat {
			return true
		}
	}
	return false
}

// difficultyOrder returns a numeric ordering for difficulty levels.
func difficultyOrder(d model.Difficulty) int {
	switch d {
	case model.DifficultyBeginner:
		return 0
	case model.DifficultyIntermediate:
		return 1
	case model.DifficultyAdvanced:
		return 2
	default:
		return 3
	}
}
