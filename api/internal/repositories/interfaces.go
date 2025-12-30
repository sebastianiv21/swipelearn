package repositories

import (
	"time"

	"github.com/google/uuid"
	"swipelearn-api/internal/models"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(id uuid.UUID, updates map[string]any) (*models.User, error)
	Delete(id uuid.UUID) error
}

// DeckRepositoryInterface defines the interface for deck repository operations
type DeckRepositoryInterface interface {
	Create(deck *models.Deck) (*models.Deck, error)
	GetByID(id uuid.UUID) (*models.Deck, error)
	GetAll() ([]*models.Deck, error)
	GetByUser(userID uuid.UUID) ([]*models.Deck, error)
	Update(id uuid.UUID, updates map[string]any) (*models.Deck, error)
	Delete(id uuid.UUID) error
	GetDeckFlashcardCount(deckID uuid.UUID) (int, error)
}

// FlashcardRepositoryInterface defines the interface for flashcard repository operations
type FlashcardRepositoryInterface interface {
	Create(card *models.Flashcard) (*models.Flashcard, error)
	GetByID(id uuid.UUID) (*models.Flashcard, error)
	GetByUser(userID uuid.UUID) ([]*models.Flashcard, error)
	Update(id uuid.UUID, updates *models.UpdateFlashcardRequest) (*models.Flashcard, error)
	Delete(id uuid.UUID) error
}

// RefreshTokenRepositoryInterface defines the interface for refresh token repository operations
type RefreshTokenRepositoryInterface interface {
	StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error
	GetValidRefreshToken(userID uuid.UUID, tokenString string) (interface{}, error)
	RevokeToken(tokenID uuid.UUID) error
	RevokeUserTokens(userID uuid.UUID) error
}
