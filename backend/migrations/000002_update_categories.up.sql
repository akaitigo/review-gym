-- Update category system to match PRD Phase 0 Step 2 definitions:
-- security, performance, design, readability, error-handling
-- Previous categories: security, performance, design, correctness, maintainability

-- Step 1: Drop existing CHECK constraints
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_category_check;
ALTER TABLE review_comments DROP CONSTRAINT IF EXISTS review_comments_category_check;
ALTER TABLE reference_reviews DROP CONSTRAINT IF EXISTS reference_reviews_category_check;

-- Step 2: Migrate existing data from old categories to new categories
UPDATE exercises SET category = 'readability' WHERE category = 'maintainability';
UPDATE exercises SET category = 'error-handling' WHERE category = 'correctness';

UPDATE review_comments SET category = 'readability' WHERE category = 'maintainability';
UPDATE review_comments SET category = 'error-handling' WHERE category = 'correctness';

UPDATE reference_reviews SET category = 'readability' WHERE category = 'maintainability';
UPDATE reference_reviews SET category = 'error-handling' WHERE category = 'correctness';

-- Step 3: Add new CHECK constraints with updated categories
ALTER TABLE exercises
    ADD CONSTRAINT exercises_category_check
    CHECK (category IN ('security', 'performance', 'design', 'readability', 'error-handling'));

ALTER TABLE review_comments
    ADD CONSTRAINT review_comments_category_check
    CHECK (category IN ('security', 'performance', 'design', 'readability', 'error-handling'));

ALTER TABLE reference_reviews
    ADD CONSTRAINT reference_reviews_category_check
    CHECK (category IN ('security', 'performance', 'design', 'readability', 'error-handling'));

-- Step 4: Add category_tags to exercises for multi-category exercises
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS category_tags JSONB NOT NULL DEFAULT '[]';

-- Step 5: Add practice tracking columns to user_profiles
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS total_exercises_completed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS consecutive_days INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS last_practice_at TIMESTAMPTZ;

-- Step 6: Add attempt tracking to scores
ALTER TABLE scores ADD COLUMN IF NOT EXISTS attempt_number INTEGER NOT NULL DEFAULT 1;
ALTER TABLE scores ADD COLUMN IF NOT EXISTS duration_seconds INTEGER;

-- Step 7: Add unique constraint for user-exercise-attempt combination
CREATE UNIQUE INDEX IF NOT EXISTS idx_scores_user_exercise_attempt
    ON scores(user_id, exercise_id, attempt_number);
