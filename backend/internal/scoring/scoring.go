// Package scoring implements the scoring engine that compares user review
// comments against reference reviews to compute precision, recall, and
// category-level scores.
package scoring

import (
	"encoding/json"
	"math"

	"github.com/akaitigo/review-gym/internal/model"
)

// lineProximityThreshold defines the maximum line distance for a match.
const lineProximityThreshold = 3

// severityWeight returns the scoring weight for a given severity.
func severityWeight(s model.Severity) float64 {
	switch s {
	case model.SeverityCritical:
		return 3.0
	case model.SeverityMajor:
		return 2.0
	case model.SeverityMinor:
		return 1.0
	case model.SeverityInfo:
		return 0.5
	default:
		return 1.0
	}
}

// Match represents a matched pair of user comment and reference review.
type Match struct {
	UserComment     model.ReviewComment   `json:"user_comment"`
	ReferenceReview model.ReferenceReview `json:"reference_review"`
	LineDelta       int                   `json:"line_delta"`
}

// Result holds the full scoring output for a review session.
type Result struct {
	PrecisionScore float64                 `json:"precision_score"`
	RecallScore    float64                 `json:"recall_score"`
	OverallScore   float64                 `json:"overall_score"`
	CategoryScores []model.CategoryScore   `json:"category_scores"`
	Matches        []Match                 `json:"matches"`
	MissedReviews  []model.ReferenceReview `json:"missed_reviews"`
	FalsePositives []model.ReviewComment   `json:"false_positives"`
}

// Compute scores a set of user comments against reference reviews.
// It uses a greedy proximity + category matching algorithm as described in ADR-003.
func Compute(comments []model.ReviewComment, references []model.ReferenceReview) Result {
	matches, matchedCommentIDs, matchedRefIDs := findMatches(comments, references)

	// Precision: fraction of user comments that matched a reference.
	precisionScore := 0.0
	if len(comments) > 0 {
		precisionScore = float64(len(matches)) / float64(len(comments)) * 100
	}

	// Recall: weighted fraction of reference reviews found.
	totalWeight := 0.0
	matchedWeight := 0.0
	for _, ref := range references {
		w := severityWeight(ref.Severity)
		totalWeight += w
		if matchedRefIDs[ref.ID] {
			matchedWeight += w
		}
	}
	recallScore := 0.0
	if totalWeight > 0 {
		recallScore = matchedWeight / totalWeight * 100
	}

	// Overall: weighted combination (40% precision, 60% recall).
	overallScore := precisionScore*0.4 + recallScore*0.6

	// Round scores to 1 decimal place.
	precisionScore = roundTo1(precisionScore)
	recallScore = roundTo1(recallScore)
	overallScore = roundTo1(overallScore)

	// Category-level scoring.
	categoryScores := computeCategoryScores(references, matchedRefIDs)

	// Collect missed reviews and false positives.
	var missed []model.ReferenceReview
	for _, ref := range references {
		if !matchedRefIDs[ref.ID] {
			missed = append(missed, ref)
		}
	}

	var falsePositives []model.ReviewComment
	for _, c := range comments {
		if !matchedCommentIDs[c.ID] {
			falsePositives = append(falsePositives, c)
		}
	}

	return Result{
		PrecisionScore: precisionScore,
		RecallScore:    recallScore,
		OverallScore:   overallScore,
		CategoryScores: categoryScores,
		Matches:        matches,
		MissedReviews:  missed,
		FalsePositives: falsePositives,
	}
}

// findMatches performs greedy matching of user comments to reference reviews.
// A match requires: same file path, same category, and line proximity <= threshold.
// Each reference can match at most one comment (closest line wins).
func findMatches(comments []model.ReviewComment, references []model.ReferenceReview) ([]Match, map[string]bool, map[string]bool) {
	matchedCommentIDs := make(map[string]bool)
	matchedRefIDs := make(map[string]bool)
	var matches []Match

	// For each reference, find the best matching comment.
	for _, ref := range references {
		bestIdx := -1
		bestDelta := lineProximityThreshold + 1

		for i, comment := range comments {
			if matchedCommentIDs[comment.ID] {
				continue
			}
			if comment.FilePath != ref.FilePath {
				continue
			}
			if comment.Category != ref.Category {
				continue
			}
			delta := abs(comment.LineNumber - ref.LineNumber)
			if delta <= lineProximityThreshold && delta < bestDelta {
				bestDelta = delta
				bestIdx = i
			}
		}

		if bestIdx >= 0 {
			matches = append(matches, Match{
				UserComment:     comments[bestIdx],
				ReferenceReview: ref,
				LineDelta:       bestDelta,
			})
			matchedCommentIDs[comments[bestIdx].ID] = true
			matchedRefIDs[ref.ID] = true
		}
	}

	return matches, matchedCommentIDs, matchedRefIDs
}

// computeCategoryScores computes per-category recall scores.
func computeCategoryScores(references []model.ReferenceReview, matchedRefIDs map[string]bool) []model.CategoryScore {
	type catAcc struct {
		totalWeight   float64
		matchedWeight float64
	}
	acc := make(map[model.Category]*catAcc)

	for _, ref := range references {
		a, ok := acc[ref.Category]
		if !ok {
			a = &catAcc{}
			acc[ref.Category] = a
		}
		w := severityWeight(ref.Severity)
		a.totalWeight += w
		if matchedRefIDs[ref.ID] {
			a.matchedWeight += w
		}
	}

	// Build scores for all categories, including those with no references.
	var scores []model.CategoryScore
	for _, cat := range model.AllCategories() {
		cs := model.CategoryScore{Category: cat}
		if a, ok := acc[cat]; ok {
			cs.MaxPoints = a.totalWeight
			cs.Earned = a.matchedWeight
			if a.totalWeight > 0 {
				cs.Score = roundTo1(a.matchedWeight / a.totalWeight * 100)
			}
		}
		scores = append(scores, cs)
	}
	return scores
}

// ToScore converts a Result into a persistable Score model.
func (r *Result) ToScore(userID, exerciseID string, attemptNumber int) (*model.Score, error) {
	categoryJSON, err := json.Marshal(r.CategoryScores)
	if err != nil {
		return nil, err
	}
	return &model.Score{
		UserID:         userID,
		ExerciseID:     exerciseID,
		PrecisionScore: r.PrecisionScore,
		RecallScore:    r.RecallScore,
		OverallScore:   r.OverallScore,
		CategoryScores: categoryJSON,
		AttemptNumber:  attemptNumber,
	}, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func roundTo1(f float64) float64 {
	return math.Round(f*10) / 10
}
