-- Add unique constraint on reference_reviews for seed idempotency.
-- Prevents duplicate reference reviews on repeated seed runs.
CREATE UNIQUE INDEX IF NOT EXISTS idx_reference_reviews_unique
    ON reference_reviews (exercise_id, file_path, line_number, content);
