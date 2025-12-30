-- Rollback authentication fields from SwipeLearn Database

-- Drop refresh_tokens table
DROP TABLE IF EXISTS refresh_tokens;

-- Remove password_hash from users table
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;