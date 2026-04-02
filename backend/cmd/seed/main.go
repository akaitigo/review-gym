package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/akaitigo/review-gym/internal/model"
	"github.com/akaitigo/review-gym/internal/seed"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://review_gym:review_gym_dev@localhost:5432/review_gym?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	exercises := seed.All()

	for _, ew := range exercises {
		exerciseID, err := insertExercise(db, &ew)
		if err != nil {
			log.Fatalf("failed to insert exercise %q: %v", ew.Exercise.Title, err)
		}

		for _, review := range ew.Reviews {
			if err := insertReferenceReview(db, exerciseID, &review); err != nil {
				log.Fatalf("failed to insert reference review for exercise %q: %v", ew.Exercise.Title, err)
			}
		}

		fmt.Printf("  Seeded: %s (%d reference reviews)\n", ew.Exercise.Title, len(ew.Reviews))
	}

	fmt.Printf("\nSuccessfully seeded %d exercises.\n", len(exercises))
}

func insertExercise(db *sql.DB, ew *seed.ExerciseWithReviews) (string, error) {
	e := ew.Exercise

	categoryTags, err := json.Marshal(e.CategoryTags)
	if err != nil {
		return "", fmt.Errorf("marshal category_tags: %w", err)
	}

	filePaths, err := json.Marshal(e.FilePaths)
	if err != nil {
		return "", fmt.Errorf("marshal file_paths: %w", err)
	}

	metadata := e.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}

	// Use a CTE to insert-or-select: if the title already exists, return
	// the existing row's ID instead of inserting a duplicate.
	var id string
	err = db.QueryRow(`
		WITH ins AS (
			INSERT INTO exercises (title, description, difficulty, category, category_tags, language, diff_content, file_paths, metadata, is_published)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (title) DO NOTHING
			RETURNING id
		)
		SELECT id FROM ins
		UNION ALL
		SELECT id FROM exercises WHERE title = $1
		LIMIT 1
	`, e.Title, e.Description, string(e.Difficulty), string(e.Category),
		categoryTags, e.Language, e.DiffContent, filePaths, metadata, e.IsPublished,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert exercise: %w", err)
	}

	return id, nil
}

func insertReferenceReview(db *sql.DB, exerciseID string, r *model.ReferenceReview) error {
	_, err := db.Exec(`
		INSERT INTO reference_reviews (exercise_id, file_path, line_number, content, category, severity, explanation)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (exercise_id, file_path, line_number, content) DO NOTHING
	`, exerciseID, r.FilePath, r.LineNumber, r.Content, string(r.Category), string(r.Severity), r.Explanation)
	if err != nil {
		return fmt.Errorf("insert reference review: %w", err)
	}

	return nil
}
