-- Add user ownership to decks table

-- Add user_id column to decks table
ALTER TABLE decks ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Create index for user_id lookup
CREATE INDEX IF NOT EXISTS idx_decks_user_id ON decks(user_id);

-- For existing decks without user_id, we need to handle this case
-- In production, you might want to assign existing decks to a system user or delete them
-- For now, we'll set user_id to NULL for existing decks (they'll need to be migrated manually)