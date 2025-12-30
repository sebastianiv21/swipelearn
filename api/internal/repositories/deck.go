package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"swipelearn-api/internal/models"
)

type DeckRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewDeckRepository(db *sql.DB, logger *logrus.Logger) *DeckRepository {
	return &DeckRepository{
		DB:     db,
		Logger: logger,
	}
}

// Create creates a new deck
func (r *DeckRepository) Create(deck *models.Deck) (*models.Deck, error) {
	query := `
		INSERT INTO decks (id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.DB.QueryRow(
		query,
		deck.ID,
		deck.Name,
		deck.Description,
	).Scan(
		&deck.ID,
		&deck.Name,
		&deck.Description,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if err != nil {
		r.Logger.WithError(err).Error("Failed to create deck in database")
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	r.Logger.WithField("deck_id", deck.ID).Info("Deck created successfully")
	return deck, nil
}

// GetByID retrieves a deck by ID
func (r *DeckRepository) GetByID(id uuid.UUID) (*models.Deck, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM decks
		WHERE id = $1
	`

	deck := &models.Deck{}
	err := r.DB.QueryRow(query, id).Scan(
		&deck.ID,
		&deck.Name,
		&deck.Description,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deck not found")
		}
		r.Logger.WithError(err).WithField("deck_id", id).Error("Failed to get deck by ID")
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	return deck, nil
}

// GetAll retrieves all decks
func (r *DeckRepository) GetAll() ([]*models.Deck, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM decks
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to get all decks")
		return nil, fmt.Errorf("failed to get decks: %w", err)
	}
	defer rows.Close()

	var decks []*models.Deck
	for rows.Next() {
		deck := &models.Deck{}
		err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&deck.Description,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			r.Logger.WithError(err).Error("Failed to scan deck row")
			return nil, fmt.Errorf("failed to scan deck: %w", err)
		}
		decks = append(decks, deck)
	}

	if err = rows.Err(); err != nil {
		r.Logger.WithError(err).Error("Error after scanning deck rows")
		return nil, fmt.Errorf("error scanning decks: %w", err)
	}

	return decks, nil
}

// Update updates a deck
func (r *DeckRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.Deck, error) {
	// Build dynamic UPDATE query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if name, ok := updates["name"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if description, ok := updates["description"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, description)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add updated_at and id
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE decks
		SET %s
		WHERE id = $%d
		RETURNING id, name, description, created_at, updated_at
	`, fmt.Sprintf("%s", setParts), argIndex)

	deck := &models.Deck{}
	err := r.DB.QueryRow(query, args...).Scan(
		&deck.ID,
		&deck.Name,
		&deck.Description,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deck not found")
		}
		r.Logger.WithError(err).WithField("deck_id", id).Error("Failed to update deck")
		return nil, fmt.Errorf("failed to update deck: %w", err)
	}

	r.Logger.WithField("deck_id", deck.ID).Info("Deck updated successfully")
	return deck, nil
}

// Delete deletes a deck by ID
func (r *DeckRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM decks WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		r.Logger.WithError(err).WithField("deck_id", id).Error("Failed to delete deck")
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deck not found")
	}

	r.Logger.WithField("deck_id", id).Info("Deck deleted successfully")
	return nil
}

// GetDeckFlashcardCount retrieves the number of flashcards in a deck
func (r *DeckRepository) GetDeckFlashcardCount(deckID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM flashcards WHERE deck_id = $1`

	var count int
	err := r.DB.QueryRow(query, deckID).Scan(&count)
	if err != nil {
		r.Logger.WithError(err).WithField("deck_id", deckID).Error("Failed to get flashcard count for deck")
		return 0, fmt.Errorf("failed to get flashcard count: %w", err)
	}

	return count, nil
}
