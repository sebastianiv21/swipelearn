DROP INDEX IF EXISTS idx_flashcards_next_review;
DROP INDEX IF EXISTS idx_flashcards_deck_id;
DROP INDEX IF EXISTS idx_flashcards_user_id;

DROP TABLE IF EXISTS flashcards;
DROP TABLE IF EXISTS decks;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
