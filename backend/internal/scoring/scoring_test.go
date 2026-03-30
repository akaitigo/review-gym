package scoring

import (
	"math"
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestCompute_PerfectScore(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
		{
			ID:         "ref-2",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Error not checked after database call",
			Category:   model.CategoryErrorHandling,
			Severity:   model.SeverityMajor,
		},
	}

	// Identical content gives perfect content similarity.
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input",
			Category:   model.CategorySecurity,
		},
		{
			ID:         "c-2",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Error not checked after database call",
			Category:   model.CategoryErrorHandling,
		},
	}

	result := Compute(comments, refs)

	if result.PrecisionScore != 100 {
		t.Errorf("precision = %.1f, want 100", result.PrecisionScore)
	}
	if result.RecallScore != 100 {
		t.Errorf("recall = %.1f, want 100", result.RecallScore)
	}
	if result.OverallScore != 100 {
		t.Errorf("overall = %.1f, want 100", result.OverallScore)
	}
	if len(result.Matches) != 2 {
		t.Errorf("matches = %d, want 2", len(result.Matches))
	}
	if len(result.MissedReviews) != 0 {
		t.Errorf("missed = %d, want 0", len(result.MissedReviews))
	}
	if len(result.FalsePositives) != 0 {
		t.Errorf("false positives = %d, want 0", len(result.FalsePositives))
	}
}

func TestCompute_NoComments(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	result := Compute(nil, refs)

	if result.PrecisionScore != 0 {
		t.Errorf("precision = %.1f, want 0", result.PrecisionScore)
	}
	if result.RecallScore != 0 {
		t.Errorf("recall = %.1f, want 0", result.RecallScore)
	}
	if len(result.MissedReviews) != 1 {
		t.Errorf("missed = %d, want 1", len(result.MissedReviews))
	}
}

func TestCompute_NoReferences(t *testing.T) {
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Random comment",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, nil)

	if result.PrecisionScore != 0 {
		t.Errorf("precision = %.1f, want 0", result.PrecisionScore)
	}
	if result.RecallScore != 0 {
		t.Errorf("recall = %.1f, want 0", result.RecallScore)
	}
	if len(result.FalsePositives) != 1 {
		t.Errorf("false positives = %d, want 1", len(result.FalsePositives))
	}
}

func TestCompute_LineProximity(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Security issue found here",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
	}

	tests := []struct {
		name        string
		commentLine int
		wantMatch   bool
	}{
		{"exact match", 10, true},
		{"within threshold +3", 13, true},
		{"within threshold -3", 7, true},
		{"outside threshold +4", 14, false},
		{"outside threshold -4", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comments := []model.ReviewComment{
				{
					ID:         "c-1",
					ExerciseID: "ex-1",
					UserID:     "user-1",
					FilePath:   "main.go",
					LineNumber: tt.commentLine,
					Content:    "Security issue found here",
					Category:   model.CategorySecurity,
				},
			}

			result := Compute(comments, refs)
			gotMatch := len(result.Matches) > 0
			if gotMatch != tt.wantMatch {
				t.Errorf("line %d: matched = %v, want %v", tt.commentLine, gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestCompute_CategoryMismatch(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Security issue",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Performance issue",
			Category:   model.CategoryPerformance, // wrong category
		},
	}

	result := Compute(comments, refs)

	if len(result.Matches) != 0 {
		t.Errorf("matches = %d, want 0 (category mismatch)", len(result.Matches))
	}
	if len(result.MissedReviews) != 1 {
		t.Errorf("missed = %d, want 1", len(result.MissedReviews))
	}
	if len(result.FalsePositives) != 1 {
		t.Errorf("false positives = %d, want 1", len(result.FalsePositives))
	}
}

func TestCompute_FilePathMismatch(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Issue",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
	}

	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "other.go", // different file
			LineNumber: 10,
			Content:    "Found issue",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	if len(result.Matches) != 0 {
		t.Errorf("matches = %d, want 0 (file path mismatch)", len(result.Matches))
	}
}

func TestCompute_SeverityWeighting(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Critical security vulnerability in authentication",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical, // weight 3.0
		},
		{
			ID:         "ref-2",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Readability issue with variable naming",
			Category:   model.CategoryReadability,
			Severity:   model.SeverityInfo, // weight 0.5
		},
	}

	// User only finds the critical issue.
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Critical security vulnerability in authentication",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	// Recall: 3.0 / (3.0 + 0.5) = 85.7%
	expectedRecall := 85.7
	if result.RecallScore != expectedRecall {
		t.Errorf("recall = %.1f, want %.1f", result.RecallScore, expectedRecall)
	}
}

func TestCompute_CategoryScores(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in query builder",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
		{
			ID:         "ref-2",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Cross-site scripting in template rendering",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
		{
			ID:         "ref-3",
			FilePath:   "main.go",
			LineNumber: 30,
			Content:    "Performance issue with N+1 query",
			Category:   model.CategoryPerformance,
			Severity:   model.SeverityMinor,
		},
	}

	// User finds only one security issue with matching content.
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in query builder",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	if len(result.CategoryScores) != 5 {
		t.Fatalf("category scores count = %d, want 5", len(result.CategoryScores))
	}

	// Check security category: 3.0 / (3.0 + 2.0) = 60%
	secScore := result.CategoryScores[0]
	if secScore.Category != model.CategorySecurity {
		t.Errorf("first category = %q, want security", secScore.Category)
	}
	if secScore.Score != 60 {
		t.Errorf("security score = %.1f, want 60", secScore.Score)
	}

	// Check performance category: 0 / 1.0 = 0%
	perfScore := result.CategoryScores[1]
	if perfScore.Category != model.CategoryPerformance {
		t.Errorf("second category = %q, want performance", perfScore.Category)
	}
	if perfScore.Score != 0 {
		t.Errorf("performance score = %.1f, want 0", perfScore.Score)
	}
}

func TestCompute_OneToOneMatching(t *testing.T) {
	// Two comments on the same line should only match one reference.
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Security vulnerability found in authentication logic",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
	}

	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Security vulnerability in authentication",
			Category:   model.CategorySecurity,
		},
		{
			ID:         "c-2",
			FilePath:   "main.go",
			LineNumber: 11,
			Content:    "Authentication logic has security vulnerability",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	if len(result.Matches) != 1 {
		t.Errorf("matches = %d, want 1 (one-to-one)", len(result.Matches))
	}
	// The closer comment (c-1 at line 10) should be matched.
	if len(result.Matches) > 0 && result.Matches[0].UserComment.ID != "c-1" {
		t.Errorf("matched comment = %q, want c-1", result.Matches[0].UserComment.ID)
	}
	if len(result.FalsePositives) != 1 {
		t.Errorf("false positives = %d, want 1", len(result.FalsePositives))
	}
}

func TestToScore(t *testing.T) {
	result := Result{
		PrecisionScore: 75.0,
		RecallScore:    50.0,
		OverallScore:   60.0,
		CategoryScores: []model.CategoryScore{
			{Category: model.CategorySecurity, Score: 100, MaxPoints: 3, Earned: 3},
		},
	}

	score, err := result.ToScore("user-1", "ex-1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.UserID != "user-1" {
		t.Errorf("user_id = %q, want user-1", score.UserID)
	}
	if score.ExerciseID != "ex-1" {
		t.Errorf("exercise_id = %q, want ex-1", score.ExerciseID)
	}
	if score.PrecisionScore != 75 {
		t.Errorf("precision = %.1f, want 75", score.PrecisionScore)
	}
	if score.AttemptNumber != 1 {
		t.Errorf("attempt = %d, want 1", score.AttemptNumber)
	}
	if err := score.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

// --- Content validation tests ---

func TestCompute_EmptyContentLowScore(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	// Empty content should not match (below similarity threshold).
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	if len(result.Matches) != 0 {
		t.Errorf("matches = %d, want 0 (empty content should not match)", len(result.Matches))
	}
	if result.PrecisionScore != 0 {
		t.Errorf("precision = %.1f, want 0 (empty content)", result.PrecisionScore)
	}
	if len(result.FalsePositives) != 1 {
		t.Errorf("false positives = %d, want 1", len(result.FalsePositives))
	}
	if len(result.MissedReviews) != 1 {
		t.Errorf("missed = %d, want 1", len(result.MissedReviews))
	}
}

func TestCompute_ShortContentPenalty(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	// Short content (< 10 chars) that still matches some words should get
	// a lower score than a full-length content with the same similarity.
	shortComment := []model.ReviewComment{
		{
			ID:         "c-short",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL input", // 9 chars, some overlap
			Category:   model.CategorySecurity,
		},
	}

	longComment := []model.ReviewComment{
		{
			ID:         "c-long",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
		},
	}

	shortResult := Compute(shortComment, refs)
	longResult := Compute(longComment, refs)

	if longResult.PrecisionScore <= shortResult.PrecisionScore {
		t.Errorf("long content precision (%.1f) should be > short content precision (%.1f)",
			longResult.PrecisionScore, shortResult.PrecisionScore)
	}
}

func TestCompute_HighContentSimilarityBoostsScore(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	// High similarity content.
	highSimComments := []model.ReviewComment{
		{
			ID:         "c-high",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
		},
	}

	// Lower similarity but still matching content.
	lowSimComments := []model.ReviewComment{
		{
			ID:         "c-low",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "There is a vulnerability in input validation logic here",
			Category:   model.CategorySecurity,
		},
	}

	highResult := Compute(highSimComments, refs)
	lowResult := Compute(lowSimComments, refs)

	if highResult.PrecisionScore <= lowResult.PrecisionScore {
		t.Errorf("high similarity precision (%.1f) should be > low similarity precision (%.1f)",
			highResult.PrecisionScore, lowResult.PrecisionScore)
	}
}

func TestCompute_CorrectPositionIrrelevantContent(t *testing.T) {
	refs := []model.ReferenceReview{
		{
			ID:         "ref-1",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection vulnerability in user input handling",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
	}

	// Position and category are correct, but content is completely unrelated.
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "The weather today is sunny and warm outside",
			Category:   model.CategorySecurity,
		},
	}

	result := Compute(comments, refs)

	// Unrelated content should not match (below similarity threshold).
	if len(result.Matches) != 0 {
		t.Errorf("matches = %d, want 0 (irrelevant content should not match)", len(result.Matches))
	}
	if result.PrecisionScore != 0 {
		t.Errorf("precision = %.1f, want 0 (irrelevant content)", result.PrecisionScore)
	}
	if result.RecallScore != 0 {
		t.Errorf("recall = %.1f, want 0 (irrelevant content)", result.RecallScore)
	}
}

func TestContentSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{"identical", "SQL injection vulnerability", "SQL injection vulnerability", 1.0},
		{"empty both", "", "", 0.0},
		{"empty one", "SQL injection", "", 0.0},
		{"no overlap", "hello world", "foo bar baz", 0.0},
		{"partial overlap", "SQL injection vulnerability", "SQL injection found", 0.5}, // 2/(3+4-2)
		{"japanese identical", "SQLインジェクションの脆弱性", "SQLインジェクションの脆弱性", 1.0},
		{"japanese paraphrase", "SQLインジェクションの脆弱性があります", "SQLインジェクション脆弱性が存在する", -1}, // > 0.0 check below
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contentSimilarity(tt.a, tt.b)
			if tt.want == -1 {
				// Special case: just check > 0
				if got <= 0 {
					t.Errorf("contentSimilarity(%q, %q) = %.3f, want > 0", tt.a, tt.b, got)
				}
			} else if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("contentSimilarity(%q, %q) = %.3f, want %.3f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMatchQuality(t *testing.T) {
	// Perfect match: line delta 0, content similarity 1.0, long content.
	perfect := Match{
		UserComment:       model.ReviewComment{Content: "SQL injection vulnerability in user input"},
		LineDelta:         0,
		ContentSimilarity: 1.0,
	}
	if q := matchQuality(&perfect); math.Abs(q-1.0) > 0.01 {
		t.Errorf("perfect match quality = %.3f, want 1.0", q)
	}

	// Short content match: should be penalised.
	short := Match{
		UserComment:       model.ReviewComment{Content: "SQL"},
		LineDelta:         0,
		ContentSimilarity: 1.0,
	}
	perfectQ := matchQuality(&perfect)
	shortQ := matchQuality(&short)
	if shortQ >= perfectQ {
		t.Errorf("short content quality (%.3f) should be < perfect quality (%.3f)", shortQ, perfectQ)
	}
}

func TestSeverityWeight(t *testing.T) {
	tests := []struct {
		severity model.Severity
		want     float64
	}{
		{model.SeverityCritical, 3.0},
		{model.SeverityMajor, 2.0},
		{model.SeverityMinor, 1.0},
		{model.SeverityInfo, 0.5},
	}

	for _, tt := range tests {
		got := severityWeight(tt.severity)
		if got != tt.want {
			t.Errorf("severityWeight(%q) = %.1f, want %.1f", tt.severity, got, tt.want)
		}
	}
}
