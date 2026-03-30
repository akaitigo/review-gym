-- Rollback initial schema
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS reference_reviews;
DROP TABLE IF EXISTS review_comments;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS user_profiles;
DROP EXTENSION IF EXISTS "uuid-ossp";
