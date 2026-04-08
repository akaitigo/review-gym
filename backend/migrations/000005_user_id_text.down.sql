-- Revert user_id columns from TEXT back to UUID and restore foreign key constraints.

-- Step 1: Drop the unique index.
DROP INDEX IF EXISTS idx_scores_user_exercise_attempt;

-- Step 2: Change column types back to UUID.
ALTER TABLE review_comments ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE scores ALTER COLUMN user_id TYPE UUID USING user_id::UUID;

-- Step 3: Restore foreign key constraints.
ALTER TABLE review_comments
    ADD CONSTRAINT review_comments_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES user_profiles(id) ON DELETE CASCADE;

ALTER TABLE scores
    ADD CONSTRAINT scores_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES user_profiles(id) ON DELETE CASCADE;

-- Step 4: Recreate the unique index.
CREATE UNIQUE INDEX IF NOT EXISTS idx_scores_user_exercise_attempt
    ON scores(user_id, exercise_id, attempt_number);
