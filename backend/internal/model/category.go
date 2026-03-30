// Package model defines core domain types for the review-gym application.
package model

import "fmt"

// Category represents a review comment category.
type Category string

const (
	CategorySecurity      Category = "security"
	CategoryPerformance   Category = "performance"
	CategoryDesign        Category = "design"
	CategoryReadability   Category = "readability"
	CategoryErrorHandling Category = "error-handling"
)

// AllCategories returns all valid categories in display order.
func AllCategories() []Category {
	return []Category{
		CategorySecurity,
		CategoryPerformance,
		CategoryDesign,
		CategoryReadability,
		CategoryErrorHandling,
	}
}

// Valid reports whether c is a recognized category.
func (c Category) Valid() bool {
	switch c {
	case CategorySecurity, CategoryPerformance, CategoryDesign,
		CategoryReadability, CategoryErrorHandling:
		return true
	default:
		return false
	}
}

// String returns the string representation of the category.
func (c Category) String() string {
	return string(c)
}

// ParseCategory converts a raw string to a Category, returning an error for unknown values.
func ParseCategory(s string) (Category, error) {
	c := Category(s)
	if !c.Valid() {
		return "", fmt.Errorf("unknown category: %q", s)
	}
	return c, nil
}
