-- Initial schema for review-gym
-- Tables: user_profiles, exercises, review_comments, reference_reviews, scores

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    display_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    github_id VARCHAR(100) UNIQUE,
    avatar_url TEXT,
    weakness_categories JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    category VARCHAR(50) NOT NULL CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability')),
    language VARCHAR(50) NOT NULL,
    diff_content TEXT NOT NULL,
    file_paths JSONB NOT NULL DEFAULT '[]',
    metadata JSONB NOT NULL DEFAULT '{}',
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE review_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    line_number INTEGER NOT NULL CHECK (line_number > 0),
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reference_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    line_number INTEGER NOT NULL CHECK (line_number > 0),
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('security', 'performance', 'design', 'correctness', 'maintainability')),
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('critical', 'major', 'minor', 'info')),
    explanation TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    precision_score NUMERIC(5,2) NOT NULL CHECK (precision_score >= 0 AND precision_score <= 100),
    recall_score NUMERIC(5,2) NOT NULL CHECK (recall_score >= 0 AND recall_score <= 100),
    overall_score NUMERIC(5,2) NOT NULL CHECK (overall_score >= 0 AND overall_score <= 100),
    category_scores JSONB NOT NULL DEFAULT '{}',
    feedback TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_exercises_category ON exercises(category);
CREATE INDEX idx_exercises_difficulty ON exercises(difficulty);
CREATE INDEX idx_exercises_published ON exercises(is_published);
CREATE INDEX idx_review_comments_exercise ON review_comments(exercise_id);
CREATE INDEX idx_review_comments_user ON review_comments(user_id);
CREATE INDEX idx_reference_reviews_exercise ON reference_reviews(exercise_id);
CREATE INDEX idx_scores_user ON scores(user_id);
CREATE INDEX idx_scores_exercise ON scores(exercise_id);
CREATE INDEX idx_scores_user_exercise ON scores(user_id, exercise_id);
