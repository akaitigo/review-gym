// Package store provides data access abstractions for the review-gym application.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/akaitigo/review-gym/internal/model"
)

// PostgresStore is a PostgreSQL implementation of all store interfaces.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a PostgresStore connected to the given database URL.
// It validates the connection with a ping before returning.
func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("ping database: %w (close error: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// Close closes the underlying database connection.
func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}

// DB returns the underlying database connection for use in migrations or seeding.
func (ps *PostgresStore) DB() *sql.DB {
	return ps.db
}

// List returns exercises matching the given filter.
func (ps *PostgresStore) List(ctx context.Context, filter ExerciseFilter) ([]model.Exercise, error) {
	query := `
		SELECT id, title, description, difficulty, category, category_tags,
		       language, diff_content, file_paths, metadata, is_published,
		       created_at, updated_at
		FROM exercises
		WHERE is_published = true`

	var args []interface{}
	argPos := 1

	if filter.Category != nil {
		query += fmt.Sprintf(` AND (category = $%d OR category_tags @> $%d::jsonb)`, argPos, argPos+1)
		args = append(args, string(*filter.Category))
		categoryJSON, err := json.Marshal([]string{string(*filter.Category)})
		if err != nil {
			return nil, fmt.Errorf("marshal category filter: %w", err)
		}
		args = append(args, string(categoryJSON))
		argPos += 2
	}

	if filter.Difficulty != nil {
		query += fmt.Sprintf(` AND difficulty = $%d`, argPos)
		args = append(args, string(*filter.Difficulty))
	}

	query += ` ORDER BY created_at ASC`

	rows, err := ps.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query exercises: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("rows.Close error: %v", closeErr)
		}
	}()

	var exercises []model.Exercise
	for rows.Next() {
		var ex model.Exercise
		var categoryTagsJSON, filePathsJSON, metadataJSON []byte

		if err := rows.Scan(
			&ex.ID, &ex.Title, &ex.Description, &ex.Difficulty, &ex.Category,
			&categoryTagsJSON, &ex.Language, &ex.DiffContent, &filePathsJSON,
			&metadataJSON, &ex.IsPublished, &ex.CreatedAt, &ex.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan exercise: %w", err)
		}

		if err := json.Unmarshal(categoryTagsJSON, &ex.CategoryTags); err != nil {
			return nil, fmt.Errorf("unmarshal category_tags: %w", err)
		}
		if err := json.Unmarshal(filePathsJSON, &ex.FilePaths); err != nil {
			return nil, fmt.Errorf("unmarshal file_paths: %w", err)
		}
		ex.Metadata = json.RawMessage(metadataJSON)

		exercises = append(exercises, ex)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate exercises: %w", err)
	}

	return exercises, nil
}

// GetByID returns a single exercise by ID.
func (ps *PostgresStore) GetByID(ctx context.Context, id string) (*model.Exercise, error) {
	query := `
		SELECT id, title, description, difficulty, category, category_tags,
		       language, diff_content, file_paths, metadata, is_published,
		       created_at, updated_at
		FROM exercises
		WHERE id = $1`

	var ex model.Exercise
	var categoryTagsJSON, filePathsJSON, metadataJSON []byte

	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&ex.ID, &ex.Title, &ex.Description, &ex.Difficulty, &ex.Category,
		&categoryTagsJSON, &ex.Language, &ex.DiffContent, &filePathsJSON,
		&metadataJSON, &ex.IsPublished, &ex.CreatedAt, &ex.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query exercise by id: %w", err)
	}

	if err := json.Unmarshal(categoryTagsJSON, &ex.CategoryTags); err != nil {
		return nil, fmt.Errorf("unmarshal category_tags: %w", err)
	}
	if err := json.Unmarshal(filePathsJSON, &ex.FilePaths); err != nil {
		return nil, fmt.Errorf("unmarshal file_paths: %w", err)
	}
	ex.Metadata = json.RawMessage(metadataJSON)

	return &ex, nil
}

// Create stores a new review comment.
func (ps *PostgresStore) Create(ctx context.Context, comment *model.ReviewComment) error {
	query := `
		INSERT INTO review_comments (exercise_id, user_id, file_path, line_number, content, category)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := ps.db.QueryRowContext(ctx, query,
		comment.ExerciseID, comment.UserID, comment.FilePath,
		comment.LineNumber, comment.Content, string(comment.Category),
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert review comment: %w", err)
	}

	return nil
}

// ListByExerciseAndUser returns comments for a given exercise and user.
func (ps *PostgresStore) ListByExerciseAndUser(ctx context.Context, exerciseID, userID string) ([]model.ReviewComment, error) {
	query := `
		SELECT id, exercise_id, user_id, file_path, line_number, content, category,
		       created_at, updated_at
		FROM review_comments
		WHERE exercise_id = $1 AND user_id = $2
		ORDER BY created_at ASC`

	rows, err := ps.db.QueryContext(ctx, query, exerciseID, userID)
	if err != nil {
		return nil, fmt.Errorf("query review comments: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("rows.Close error: %v", closeErr)
		}
	}()

	var comments []model.ReviewComment
	for rows.Next() {
		var c model.ReviewComment
		if err := rows.Scan(
			&c.ID, &c.ExerciseID, &c.UserID, &c.FilePath, &c.LineNumber,
			&c.Content, &c.Category, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan review comment: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate review comments: %w", err)
	}

	return comments, nil
}

// ListByExercise returns all reference reviews for an exercise.
func (ps *PostgresStore) ListByExercise(ctx context.Context, exerciseID string) ([]model.ReferenceReview, error) {
	query := `
		SELECT id, exercise_id, file_path, line_number, content, category,
		       severity, explanation, created_at
		FROM reference_reviews
		WHERE exercise_id = $1
		ORDER BY created_at ASC`

	rows, err := ps.db.QueryContext(ctx, query, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("query reference reviews: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("rows.Close error: %v", closeErr)
		}
	}()

	var reviews []model.ReferenceReview
	for rows.Next() {
		var r model.ReferenceReview
		if err := rows.Scan(
			&r.ID, &r.ExerciseID, &r.FilePath, &r.LineNumber, &r.Content,
			&r.Category, &r.Severity, &r.Explanation, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reference review: %w", err)
		}
		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reference reviews: %w", err)
	}

	return reviews, nil
}

// SaveScore persists a scoring result.
// It uses a transaction with FOR UPDATE lock to atomically compute the next
// attempt number, preventing race conditions from concurrent requests.
func (ps *PostgresStore) SaveScore(ctx context.Context, score *model.Score) error {
	categoryScores := score.CategoryScores
	if categoryScores == nil {
		categoryScores = json.RawMessage(`{}`)
	}

	var feedback sql.NullString
	if score.Feedback != "" {
		feedback = sql.NullString{String: score.Feedback, Valid: true}
	}

	var durationSeconds sql.NullInt32
	if score.DurationSeconds > 0 {
		durationSeconds = sql.NullInt32{Int32: int32(score.DurationSeconds), Valid: true}
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("tx.Rollback error: %v", rbErr)
		}
	}()

	// Lock existing score rows for this user+exercise to prevent concurrent
	// attempt number collisions. FOR UPDATE ensures serialized access.
	var maxAttempt int
	err = tx.QueryRow(`
		SELECT COALESCE(MAX(attempt_number), 0)
		FROM scores
		WHERE user_id = $1 AND exercise_id = $2
		FOR UPDATE`, score.UserID, score.ExerciseID).Scan(&maxAttempt)
	if err != nil {
		return fmt.Errorf("lock and read max attempt: %w", err)
	}

	nextAttempt := maxAttempt + 1

	insertQuery := `
		INSERT INTO scores (user_id, exercise_id, precision_score, recall_score,
		                    overall_score, category_scores, feedback, attempt_number,
		                    duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, attempt_number, created_at`

	err = tx.QueryRow(insertQuery,
		score.UserID, score.ExerciseID, score.PrecisionScore, score.RecallScore,
		score.OverallScore, string(categoryScores), feedback, nextAttempt, durationSeconds,
	).Scan(&score.ID, &score.AttemptNumber, &score.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert score: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetScoresByExerciseAndUser returns all scores for a given exercise and user,
// ordered by attempt number.
func (ps *PostgresStore) GetScoresByExerciseAndUser(ctx context.Context, exerciseID, userID string) ([]model.Score, error) {
	query := `
		SELECT id, user_id, exercise_id, precision_score, recall_score,
		       overall_score, category_scores, COALESCE(feedback, ''),
		       attempt_number, COALESCE(duration_seconds, 0), created_at
		FROM scores
		WHERE exercise_id = $1 AND user_id = $2
		ORDER BY attempt_number ASC`

	rows, err := ps.db.QueryContext(ctx, query, exerciseID, userID)
	if err != nil {
		return nil, fmt.Errorf("query scores by exercise and user: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("rows.Close error: %v", closeErr)
		}
	}()

	var scores []model.Score
	for rows.Next() {
		var s model.Score
		var categoryScoresJSON []byte
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ExerciseID, &s.PrecisionScore, &s.RecallScore,
			&s.OverallScore, &categoryScoresJSON, &s.Feedback,
			&s.AttemptNumber, &s.DurationSeconds, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan score: %w", err)
		}
		s.CategoryScores = json.RawMessage(categoryScoresJSON)
		scores = append(scores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scores: %w", err)
	}

	return scores, nil
}

// GetScoresByUser returns all scores for a given user.
func (ps *PostgresStore) GetScoresByUser(ctx context.Context, userID string) ([]model.Score, error) {
	query := `
		SELECT id, user_id, exercise_id, precision_score, recall_score,
		       overall_score, category_scores, COALESCE(feedback, ''),
		       attempt_number, COALESCE(duration_seconds, 0), created_at
		FROM scores
		WHERE user_id = $1
		ORDER BY created_at ASC`

	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query scores by user: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("rows.Close error: %v", closeErr)
		}
	}()

	var scores []model.Score
	for rows.Next() {
		var s model.Score
		var categoryScoresJSON []byte
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ExerciseID, &s.PrecisionScore, &s.RecallScore,
			&s.OverallScore, &categoryScoresJSON, &s.Feedback,
			&s.AttemptNumber, &s.DurationSeconds, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan score: %w", err)
		}
		s.CategoryScores = json.RawMessage(categoryScoresJSON)
		scores = append(scores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scores: %w", err)
	}

	return scores, nil
}

// CountCompletedExercises returns the number of distinct exercises scored by a user.
func (ps *PostgresStore) CountCompletedExercises(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(DISTINCT exercise_id) FROM scores WHERE user_id = $1`

	var count int
	if err := ps.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count completed exercises: %w", err)
	}

	return count, nil
}
