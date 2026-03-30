package scoring

import (
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
			Content:    "SQL injection",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
		{
			ID:         "ref-2",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Error not checked",
			Category:   model.CategoryErrorHandling,
			Severity:   model.SeverityMajor,
		},
	}

	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "SQL injection found",
			Category:   model.CategorySecurity,
		},
		{
			ID:         "c-2",
			ExerciseID: "ex-1",
			UserID:     "user-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Error not handled",
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
			Content:    "Issue here",
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
					Content:    "Found issue",
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
			Content:    "Critical issue",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical, // weight 3.0
		},
		{
			ID:         "ref-2",
			ExerciseID: "ex-1",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Info issue",
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
			Content:    "Found critical",
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
			Content:    "Security issue 1",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityCritical,
		},
		{
			ID:         "ref-2",
			FilePath:   "main.go",
			LineNumber: 20,
			Content:    "Security issue 2",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
		{
			ID:         "ref-3",
			FilePath:   "main.go",
			LineNumber: 30,
			Content:    "Performance issue",
			Category:   model.CategoryPerformance,
			Severity:   model.SeverityMinor,
		},
	}

	// User finds only one security issue.
	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Found it",
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
			Content:    "Issue",
			Category:   model.CategorySecurity,
			Severity:   model.SeverityMajor,
		},
	}

	comments := []model.ReviewComment{
		{
			ID:         "c-1",
			FilePath:   "main.go",
			LineNumber: 10,
			Content:    "Comment 1",
			Category:   model.CategorySecurity,
		},
		{
			ID:         "c-2",
			FilePath:   "main.go",
			LineNumber: 11,
			Content:    "Comment 2",
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
