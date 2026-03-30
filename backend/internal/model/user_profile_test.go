package model_test

import (
	"strings"
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestUserProfileValidate(t *testing.T) {
	valid := model.UserProfile{
		DisplayName:        "testuser",
		WeaknessCategories: []model.Category{model.CategorySecurity},
	}

	t.Run("valid profile passes", func(t *testing.T) {
		u := valid
		if err := u.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty display_name fails", func(t *testing.T) {
		u := valid
		u.DisplayName = ""
		err := u.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "display_name")
	})

	t.Run("display_name too long fails", func(t *testing.T) {
		u := valid
		u.DisplayName = strings.Repeat("a", 101)
		err := u.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "display_name")
	})

	t.Run("invalid weakness category fails", func(t *testing.T) {
		u := valid
		u.WeaknessCategories = []model.Category{"invalid"}
		err := u.Validate()
		if err == nil {
			t.Error("expected error")
		}
		assertValidationField(t, err, "weakness_categories")
	})
}
