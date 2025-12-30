-- Remove user ownership from decks table

-- Remove index
DROP INDEX IF EXISTS idx_decks_user_id;

-- Remove user_id column
ALTER TABLE decks DROP COLUMN IF EXISTS user_id;