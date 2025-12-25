package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"swipelearn-api/internal/utils"
)

type Database struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(logger *logrus.Logger) (*Database, error) {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set, add it to your .env file")
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(utils.GetEnvAsInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(utils.GetEnvAsInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(utils.GetEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), utils.GetEnvAsDuration("DB_CONNECT_TIMEOUT", 5*time.Second))
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully!")

	return &Database{
		DB:     db,
		Logger: logger,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	d.Logger.Info("Closing database connection")
	return d.DB.Close()
}

// Health checks database connectivity
func (d *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), utils.GetEnvAsDuration("DB_HEALTH_CHECK_TIMEOUT", 2*time.Second))
	defer cancel()

	if err := d.DB.PingContext(ctx); err != nil {
		d.Logger.WithError(err).Error("Database health check failed")
		return err
	}

	d.Logger.Info("Database health check passed")
	return nil
}
