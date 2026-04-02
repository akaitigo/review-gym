-- Add unique constraint on exercises.title for idempotent seeding.
CREATE UNIQUE INDEX IF NOT EXISTS idx_exercises_title_unique ON exercises(title);
