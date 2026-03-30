-- Rollback category updates

-- Drop new index
DROP INDEX IF EXISTS idx_scores_user_exercise_attempt;

-- Remove added columns
ALTER TABLE scores DROP COLUMN IF EXISTS duration_seconds;
ALTER TABLE scores DROP COLUMN IF EXISTS attempt_number;

ALTER TABLE user_profiles DROP COLUMN IF EXISTS last_practice_at;
ALTER TABLE user_profiles DROP COLUMN IF EXISTS consecutive_days;
ALTER TABLE user_profiles DROP COLUMN IF EXISTS total_exercises_completed;

ALTER TABLE exercises DROP COLUMN IF EXISTS category_tags;

-- Drop new CHECK constraints
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_category_check;
ALTER TABLE review_comments DROP CONSTRAINT IF EXISTS review_comments_category_check;
ALTER TABLE reference_reviews DROP CONSTRAINT IF EXISTS reference_reviews_category_check;

-- Migrate data back to original categories
UPDATE exercises SET category = 'maintainability' WHERE category = 'readability';
UPDATE exercises SET category = 'correctness' WHERE category = 'error-handling';

UPDATE review_comments SET category = 'maintainability' WHERE category = 'readability';
UPDATE review_comments SET category = 'correctness' WHERE category = 'error-handling';

UPDATE reference_reviews SET category = 'maintainability' WHERE category = 'readability';
UPDATE reference_reviews SET category = 'correctness' WHERE category = 'error-handling';

-- Restore original CHECK constraints
ALTER TABLE exercises
    ADD CONSTRAINT exercises_category_check
    CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability'));

ALTER TABLE review_comments
    ADD CONSTRAINT review_comments_category_check
    CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability'));

ALTER TABLE reference_reviews
    ADD CONSTRAINT reference_reviews_category_check
    CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability'));
