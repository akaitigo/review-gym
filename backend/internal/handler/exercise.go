package handler

import (
	"net/http"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/store"
)

// exerciseListItem is the response shape for the exercise list endpoint.
// It omits the full diff_content to reduce payload size.
type exerciseListItem struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Difficulty   model.Difficulty `json:"difficulty"`
	Category     model.Category   `json:"category"`
	CategoryTags []model.Category `json:"category_tags"`
	Language     string           `json:"language"`
	FilePaths    []string         `json:"file_paths"`
	IsPublished  bool             `json:"is_published"`
}

// ListExercises handles GET /api/exercises.
// Supports optional query parameters: category, difficulty.
func (h *Handler) ListExercises(w http.ResponseWriter, r *http.Request) {
	filter := store.ExerciseFilter{}

	if cat := r.URL.Query().Get("category"); cat != "" {
		c, err := model.ParseCategory(cat)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid category: must be one of security, performance, design, readability, error-handling")
			return
		}
		filter.Category = &c
	}

	if diff := r.URL.Query().Get("difficulty"); diff != "" {
		d, err := model.ParseDifficulty(diff)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid difficulty: must be one of beginner, intermediate, advanced")
			return
		}
		filter.Difficulty = &d
	}

	exercises, err := h.Exercises.List(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list exercises")
		return
	}

	items := make([]exerciseListItem, 0, len(exercises))
	for _, ex := range exercises {
		tags := ex.CategoryTags
		if tags == nil {
			tags = []model.Category{}
		}
		fps := ex.FilePaths
		if fps == nil {
			fps = []string{}
		}
		items = append(items, exerciseListItem{
			ID:           ex.ID,
			Title:        ex.Title,
			Description:  ex.Description,
			Difficulty:   ex.Difficulty,
			Category:     ex.Category,
			CategoryTags: tags,
			Language:     ex.Language,
			FilePaths:    fps,
			IsPublished:  ex.IsPublished,
		})
	}

	writeJSON(w, http.StatusOK, items)
}

// GetExercise handles GET /api/exercises/{id}.
func (h *Handler) GetExercise(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing exercise id")
		return
	}

	exercise, err := h.Exercises.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get exercise")
		return
	}
	if exercise == nil || !exercise.IsPublished {
		writeError(w, http.StatusNotFound, "exercise not found")
		return
	}

	writeJSON(w, http.StatusOK, exercise)
}
