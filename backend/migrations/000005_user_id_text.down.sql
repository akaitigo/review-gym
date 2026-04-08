-- Revert user_id columns from TEXT back to UUID and restore foreign key constraints.
--
-- WARNING: This migration deletes rows with non-UUID user_id values that cannot
-- be converted back to UUID. Run a backup before applying in production.

-- Step 1: Drop the unique index.
DROP INDEX IF EXISTS idx_scores_user_exercise_attempt;

-- Step 2: Clean up rows with non-UUID user_id values that would fail the cast.
-- UUID format: 8-4-4-4-12 hex characters (e.g. 00000000-0000-0000-0000-000000000000).
DELETE FROM review_comments
    WHERE user_id !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';
DELETE FROM scores
    WHERE user_id !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';

-- Step 3: Change column types back to UUID.
ALTER TABLE review_comments ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE scores ALTER COLUMN user_id TYPE UUID USING user_id::UUID;

-- Step 4: Restore foreign key constraints.
ALTER TABLE review_comments
    ADD CONSTRAINT review_comments_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES user_profiles(id) ON DELETE CASCADE;

ALTER TABLE scores
    ADD CONSTRAINT scores_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES user_profiles(id) ON DELETE CASCADE;

-- Step 5: Recreate the unique index.
CREATE UNIQUE INDEX IF NOT EXISTS idx_scores_user_exercise_attempt
    ON scores(user_id, exercise_id, attempt_number);
