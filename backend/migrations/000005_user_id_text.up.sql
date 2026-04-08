-- Change user_id columns from UUID to TEXT for review_comments and scores.
-- The current schema uses UUID with a foreign key to user_profiles, but
-- the application has no authentication yet and passes arbitrary strings
-- (e.g. "anonymous") as user IDs. This causes all API writes to fail
-- when running against PostgreSQL.
--
-- Step 1: Drop foreign key constraints referencing user_profiles.
ALTER TABLE review_comments DROP CONSTRAINT IF EXISTS review_comments_user_id_fkey;
ALTER TABLE scores DROP CONSTRAINT IF EXISTS scores_user_id_fkey;

-- Step 2: Drop the unique index on scores that includes user_id (UUID).
DROP INDEX IF EXISTS idx_scores_user_exercise_attempt;

-- Step 3: Change column types from UUID to TEXT.
ALTER TABLE review_comments ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE scores ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;

-- Step 4: Recreate the unique index with TEXT user_id.
CREATE UNIQUE INDEX IF NOT EXISTS idx_scores_user_exercise_attempt
    ON scores(user_id, exercise_id, attempt_number);
