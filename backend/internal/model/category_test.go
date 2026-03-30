package model_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestCategoryValid(t *testing.T) {
	tests := []struct {
		category model.Category
		want     bool
	}{
		{model.CategorySecurity, true},
		{model.CategoryPerformance, true},
		{model.CategoryDesign, true},
		{model.CategoryReadability, true},
		{model.CategoryErrorHandling, true},
		{"unknown", false},
		{"", false},
		{"correctness", false},
		{"maintainability", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			if got := tt.category.Valid(); got != tt.want {
				t.Errorf("Category(%q).Valid() = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestAllCategories(t *testing.T) {
	cats := model.AllCategories()
	if len(cats) != 5 {
		t.Errorf("AllCategories() returned %d categories, want 5", len(cats))
	}

	expected := []model.Category{
		model.CategorySecurity,
		model.CategoryPerformance,
		model.CategoryDesign,
		model.CategoryReadability,
		model.CategoryErrorHandling,
	}

	for i, cat := range expected {
		if cats[i] != cat {
			t.Errorf("AllCategories()[%d] = %q, want %q", i, cats[i], cat)
		}
	}
}

func TestParseCategoryValid(t *testing.T) {
	for _, raw := range []string{"security", "performance", "design", "readability", "error-handling"} {
		cat, err := model.ParseCategory(raw)
		if err != nil {
			t.Errorf("ParseCategory(%q) returned unexpected error: %v", raw, err)
		}
		if cat.String() != raw {
			t.Errorf("ParseCategory(%q).String() = %q", raw, cat.String())
		}
	}
}

func TestParseCategoryInvalid(t *testing.T) {
	_, err := model.ParseCategory("invalid")
	if err == nil {
		t.Error("ParseCategory(\"invalid\") expected error, got nil")
	}
}
