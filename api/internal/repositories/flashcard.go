package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"swipelearn-api/internal/models"
)

type FlashcardRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewFlashcardRepository(db *sql.DB, logger *logrus.Logger) *FlashcardRepository {
	return &FlashcardRepository{
		DB:     db,
		Logger: logger,
	}
}

// Create inserts a new flashcard
func (r *FlashcardRepository) Create(card *models.Flashcard) (*models.Flashcard, error) {
	query := `
        INSERT INTO flashcards (id, user_id, deck_id, front, back, difficulty, interval, ease_factor, review_count, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
        RETURNING id, user_id, deck_id, front, back, difficulty, interval, ease_factor, review_count, last_review, next_review, created_at, updated_at
    `

	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	err := r.DB.QueryRow(
		query,
		card.ID, card.UserID, card.DeckID, card.Front, card.Back,
		card.Difficulty, card.Interval, card.EaseFactor, card.ReviewCount,
	).Scan(
		&card.ID, &card.UserID, &card.DeckID, &card.Front, &card.Back,
		&card.Difficulty, &card.Interval, &card.EaseFactor, &card.ReviewCount,
		&card.LastReview, &card.NextReview, &card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		r.Logger.WithError(err).Error("Failed to create flashcard")
		return nil, fmt.Errorf("failed to create flashcard: %w", err)
	}

	r.Logger.WithFields(logrus.Fields{
		"flashcard_id": card.ID,
		"user_id":      card.UserID,
	}).Info("Flashcard created successfully")

	return card, nil
}

// GetByID retrieves a flashcard by ID
func (r *FlashcardRepository) GetByID(id uuid.UUID) (*models.Flashcard, error) {
	query := `
        SELECT id, user_id, deck_id, front, back, difficulty, interval, ease_factor, review_count,
               last_review, next_review, created_at, updated_at
        FROM flashcards
        WHERE id = $1
    `

	var card models.Flashcard
	err := r.DB.QueryRow(query, id).Scan(
		&card.ID, &card.UserID, &card.DeckID, &card.Front, &card.Back,
		&card.Difficulty, &card.Interval, &card.EaseFactor, &card.ReviewCount,
		&card.LastReview, &card.NextReview, &card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flashcard not found")
		}
		r.Logger.WithError(err).WithField("flashcard_id", id).Error("Failed to get flashcard")
		return nil, fmt.Errorf("failed to get flashcard: %w", err)
	}

	return &card, nil
}

// GetByUser retrieves all flashcards for a user
func (r *FlashcardRepository) GetByUser(userID uuid.UUID) ([]*models.Flashcard, error) {
	query := `
        SELECT id, user_id, deck_id, front, back, difficulty, interval, ease_factor, review_count,
               last_review, next_review, created_at, updated_at
        FROM flashcards
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		r.Logger.WithError(err).WithField("user_id", userID).Error("Failed to get flashcards for user")
		return nil, fmt.Errorf("failed to get flashcards: %w", err)
	}
	defer rows.Close()

	var flashcards []*models.Flashcard
	for rows.Next() {
		var card models.Flashcard
		err := rows.Scan(
			&card.ID, &card.UserID, &card.DeckID, &card.Front, &card.Back,
			&card.Difficulty, &card.Interval, &card.EaseFactor, &card.ReviewCount,
			&card.LastReview, &card.NextReview, &card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			r.Logger.WithError(err).Error("Failed to scan flashcard")
			return nil, fmt.Errorf("failed to scan flashcard: %w", err)
		}
		flashcards = append(flashcards, &card)
	}

	if err = rows.Err(); err != nil {
		r.Logger.WithError(err).Error("Error iterating flashcard rows")
		return nil, fmt.Errorf("failed to iterate flashcards: %w", err)
	}

	r.Logger.WithFields(logrus.Fields{
		"user_id":         userID,
		"flashcard_count": len(flashcards),
	}).Info("Retrieved flashcards for user")

	return flashcards, nil
}

// Update with safer named parameter approach
func (r *FlashcardRepository) Update(id uuid.UUID, updates *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	// Start with existing card
	card, err := r.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// Update fields individually
	if updates.Front != nil {
		card.Front = *updates.Front
	}
	if updates.Back != nil {
		card.Back = *updates.Back
	}
	if updates.Difficulty != nil {
		card.Difficulty = *updates.Difficulty
	}
	if updates.Interval != nil {
		card.Interval = *updates.Interval
	}
	if updates.EaseFactor != nil {
		card.EaseFactor = *updates.EaseFactor
	}
	if updates.ReviewCount != nil {
		card.ReviewCount = *updates.ReviewCount
	}
	if updates.LastReview != nil {
		card.LastReview = updates.LastReview
	}
	if updates.NextReview != nil {
		card.NextReview = updates.NextReview
	}

	// Comprehensive update query for SM-2 algorithm
	query := `
        UPDATE flashcards
        SET front = $2, back = $3, difficulty = $4, interval = $5, 
            ease_factor = $6, review_count = $7, last_review = $8, 
            next_review = $9, updated_at = NOW()
        WHERE id = $1
        RETURNING id, user_id, deck_id, front, back, difficulty, interval, ease_factor, review_count,
                  last_review, next_review, created_at, updated_at
    `

	err = r.DB.QueryRow(
		query,
		id, card.Front, card.Back, card.Difficulty, card.Interval,
		card.EaseFactor, card.ReviewCount, card.LastReview, card.NextReview,
	).Scan(
		&card.ID, &card.UserID, &card.DeckID, &card.Front, &card.Back,
		&card.Difficulty, &card.Interval, &card.EaseFactor, &card.ReviewCount,
		&card.LastReview, &card.NextReview, &card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		r.Logger.WithError(err).WithField("flashcard_id", id).Error("Failed to update flashcard")
		return nil, fmt.Errorf("failed to update flashcard: %w", err)
	}

	return card, nil
}

// Delete removes a flashcard
func (r *FlashcardRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM flashcards WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		r.Logger.WithError(err).WithField("flashcard_id", id).Error("Failed to delete flashcard")
		return fmt.Errorf("failed to delete flashcard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("flashcard not found")
	}

	r.Logger.WithField("flashcard_id", id).Info("Flashcard deleted successfully")
	return nil
}
