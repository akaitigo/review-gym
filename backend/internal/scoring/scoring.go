// Package scoring implements the scoring engine that compares user review
// comments against reference reviews to compute precision, recall, and
// category-level scores.
package scoring

import (
	"encoding/json"
	"math"
	"strings"
	"unicode"

	"github.com/akaitigo/review-gym/internal/model"
)

// lineProximityThreshold defines the maximum line distance for a match.
const lineProximityThreshold = 3

// contentSimilarityThreshold defines the minimum Jaccard similarity for content
// matching. Comments below this threshold are not considered a match even if
// file path, category, and line proximity all pass.
const contentSimilarityThreshold = 0.2

// shortContentThreshold defines the minimum character count for content to avoid
// a scoring penalty. Content shorter than this is penalised.
const shortContentThreshold = 10

// shortContentPenalty is the multiplier applied to content similarity when the
// user comment content is shorter than shortContentThreshold.
const shortContentPenalty = 0.5

// positionWeight and contentWeight control the blend of position-based and
// content-based factors in the overall match quality.
const (
	positionWeight = 0.5
	contentWeight  = 0.5
)

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
	UserComment       model.ReviewComment   `json:"user_comment"`
	ReferenceReview   model.ReferenceReview `json:"reference_review"`
	LineDelta         int                   `json:"line_delta"`
	ContentSimilarity float64               `json:"content_similarity"`
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

	// Precision: weighted fraction of user comments that matched a reference.
	// Each match contributes its quality score (position proximity + content
	// similarity) rather than counting as a simple binary hit.
	precisionScore := 0.0
	if len(comments) > 0 {
		totalQuality := 0.0
		for i := range matches {
			totalQuality += matchQuality(&matches[i])
		}
		precisionScore = totalQuality / float64(len(comments)) * 100
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
// A match requires: same file path, same category, line proximity <= threshold,
// and content similarity >= contentSimilarityThreshold.
// Each reference can match at most one comment (best combined score wins).
func findMatches(comments []model.ReviewComment, references []model.ReferenceReview) ([]Match, map[string]bool, map[string]bool) {
	matchedCommentIDs := make(map[string]bool)
	matchedRefIDs := make(map[string]bool)
	var matches []Match

	// For each reference, find the best matching comment.
	for _, ref := range references {
		bestIdx := -1
		bestDelta := lineProximityThreshold + 1
		bestSimilarity := 0.0

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
			if delta > lineProximityThreshold {
				continue
			}
			sim := contentSimilarity(comment.Content, ref.Content)
			if sim < contentSimilarityThreshold {
				continue
			}
			// Prefer closer line, then higher similarity.
			if delta < bestDelta || (delta == bestDelta && sim > bestSimilarity) {
				bestDelta = delta
				bestSimilarity = sim
				bestIdx = i
			}
		}

		if bestIdx >= 0 {
			sim := contentSimilarity(comments[bestIdx].Content, ref.Content)
			matches = append(matches, Match{
				UserComment:       comments[bestIdx],
				ReferenceReview:   ref,
				LineDelta:         bestDelta,
				ContentSimilarity: sim,
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

// tokenize splits text into a set of normalised tokens.
// For CJK characters (Japanese/Chinese/Korean), each character becomes a separate token (unigram).
// For Latin text, words are split on whitespace/punctuation as before.
func tokenize(text string) map[string]bool {
	tokens := make(map[string]bool)
	lower := strings.ToLower(text)

	// Split into words first
	words := strings.FieldsFunc(lower, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	for _, word := range words {
		// Check if the word contains CJK characters
		hasCJK := false
		for _, r := range word {
			if isCJK(r) {
				hasCJK = true
				break
			}
		}

		if hasCJK {
			// CJK: character-level unigrams + bigrams for better matching
			runes := []rune(word)
			for i, r := range runes {
				tokens[string(r)] = true
				if i+1 < len(runes) {
					tokens[string(runes[i:i+2])] = true // bigram
				}
			}
		} else {
			tokens[word] = true
		}
	}
	return tokens
}

// isCJK returns true if the rune is a CJK Unified Ideograph or Hiragana/Katakana.
func isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) || // Katakana
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Extension A
		(r >= 0xF900 && r <= 0xFAFF) // CJK Compatibility Ideographs
}

// contentSimilarity computes the Jaccard similarity coefficient between the
// word-token sets of two strings. Returns a value in [0.0, 1.0].
func contentSimilarity(a, b string) float64 {
	setA := tokenize(a)
	setB := tokenize(b)
	if len(setA) == 0 && len(setB) == 0 {
		return 0.0
	}

	intersection := 0
	for token := range setA {
		if setB[token] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0.0
	}
	return float64(intersection) / float64(union)
}

// matchQuality returns a quality score in [0.0, 1.0] for a given match.
// It blends position proximity (closer line = higher) and content similarity,
// applying a penalty for very short user comments.
func matchQuality(m *Match) float64 {
	// Position component: 1.0 for exact line, decaying to ~0.25 at threshold.
	posScore := 1.0 - float64(m.LineDelta)/float64(lineProximityThreshold+1)

	// Content component: direct similarity, with short-content penalty.
	contScore := m.ContentSimilarity
	if len(strings.TrimSpace(m.UserComment.Content)) < shortContentThreshold {
		contScore *= shortContentPenalty
	}

	return positionWeight*posScore + contentWeight*contScore
}
