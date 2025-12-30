package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// TestDatabase holds the test database connection
type TestDatabase struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
}

// SetupTestDatabase creates a test database using environment variables or default to PostgreSQL
func SetupTestDatabase(t *testing.T) *TestDatabase {
	logger := TestLogger()

	// Try to use PostgreSQL if available
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Default to PostgreSQL with test database
		databaseURL = "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"
	}

	var db *sqlx.DB
	var err error

	// Try to connect to PostgreSQL
	db, err = sqlx.Connect("postgres", databaseURL)
	if err != nil {
		// If PostgreSQL is not available, skip tests that require database
		t.Skipf("Skipping database tests: PostgreSQL not available at %s", databaseURL)
		return nil
	}

	// Test the connection
	err = db.Ping()
	require.NoError(t, err, "Could not connect to database")

	// Set connection pool settings for tests
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	// Clean up function
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	})

	return &TestDatabase{
		DB:     db,
		Logger: logger,
	}
}

// RunMigrations runs all migrations on the test database
func (td *TestDatabase) RunMigrations(t *testing.T) {
	migrations := []string{
		// Create UUID extension
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Decks table
		`CREATE TABLE IF NOT EXISTS decks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Flashcards table
		`CREATE TABLE IF NOT EXISTS flashcards (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
			front TEXT NOT NULL,
			back TEXT NOT NULL,
			difficulty FLOAT DEFAULT 2.5,
			interval INTEGER DEFAULT 1,
			ease_factor FLOAT DEFAULT 2.5,
			review_count INTEGER DEFAULT 0,
			last_review TIMESTAMP WITH TIME ZONE,
			next_review TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Refresh tokens table
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_decks_user_id ON decks(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_flashcards_user_id ON flashcards(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_flashcards_deck_id ON flashcards(deck_id);`,
		`CREATE INDEX IF NOT EXISTS idx_flashcards_next_review ON flashcards(next_review);`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);`,
	}

	for _, migration := range migrations {
		_, err := td.DB.Exec(migration)
		require.NoError(t, err, "Failed to run migration: %s", migration)
	}

	// Clear all tables to ensure clean state
	td.CleanupDatabase(t)
}

// CleanupDatabase removes all data from database tables
func (td *TestDatabase) CleanupDatabase(t *testing.T) {
	tables := []string{"refresh_tokens", "flashcards", "decks", "users"}

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("DELETE FROM %s;", table))
		require.NoError(t, err, "Failed to cleanup table: %s", table)
	}
}

// TruncateTables truncates all tables (faster than DELETE for large datasets)
func (td *TestDatabase) TruncateTables(t *testing.T) {
	tables := []string{"refresh_tokens", "flashcards", "decks", "users"}

	// Disable foreign key constraints temporarily
	_, err := td.DB.Exec("SET session_replication_role = replica;")
	require.NoError(t, err)
	defer func() {
		_, err := td.DB.Exec("SET session_replication_role = DEFAULT;")
		require.NoError(t, err)
	}()

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
		require.NoError(t, err, "Failed to truncate table: %s", table)
	}
}

// Close closes the database connection
func (td *TestDatabase) Close() error {
	if td.DB != nil {
		return td.DB.Close()
	}
	return nil
}
