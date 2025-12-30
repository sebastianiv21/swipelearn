package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"swipelearn-api/internal/models"
	"swipelearn-api/internal/repositories"
)

type DeckService struct {
	deckRepo *repositories.DeckRepository
	Logger   *logrus.Logger
}

func NewDeckService(repo *repositories.DeckRepository, logger *logrus.Logger) *DeckService {
	return &DeckService{
		deckRepo: repo,
		Logger:   logger,
	}
}

// Create creates a new deck with business logic validation
func (s *DeckService) Create(req *models.CreateDeckRequest, userID uuid.UUID) (*models.Deck, error) {
	deck := &models.Deck{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	savedDeck, err := s.deckRepo.Create(deck)
	if err != nil {
		s.Logger.WithError(err).Error("Service failed to create deck")
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"deck_id": savedDeck.ID,
		"name":    savedDeck.Name,
	}).Info("Deck created successfully")

	return savedDeck, nil
}

// GetByID retrieves a deck by ID
func (s *DeckService) GetByID(id uuid.UUID) (*models.Deck, error) {
	deck, err := s.deckRepo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("deck_id", id).Error("Service failed to get deck")
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	return deck, nil
}

// GetByIDWithOwnership retrieves a deck by ID with user ownership validation
func (s *DeckService) GetByIDWithOwnership(id uuid.UUID, userID uuid.UUID) (*models.Deck, error) {
	deck, err := s.deckRepo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("deck_id", id).Error("Service failed to get deck")
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	// Check if the deck belongs to the user
	if deck.UserID != userID {
		s.Logger.WithFields(logrus.Fields{
			"deck_id":  id,
			"user_id":  userID,
			"owner_id": deck.UserID,
		}).Warn("Unauthorized attempt to access deck")
		return nil, fmt.Errorf("unauthorized: deck does not belong to user")
	}

	return deck, nil
}

// GetAll retrieves all decks (admin only)
func (s *DeckService) GetAll() ([]*models.Deck, error) {
	decks, err := s.deckRepo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Service failed to get all decks")
		return nil, fmt.Errorf("failed to get decks: %w", err)
	}

	s.Logger.WithField("deck_count", len(decks)).Info("Retrieved all decks")
	return decks, nil
}

// GetByUser retrieves all decks for a user
func (s *DeckService) GetByUser(userID uuid.UUID) ([]*models.Deck, error) {
	decks, err := s.deckRepo.GetByUser(userID)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", userID).Error("Service failed to get decks for user")
		return nil, fmt.Errorf("failed to get decks: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"deck_count": len(decks),
	}).Info("Retrieved decks for user")
	return decks, nil
}

// Update updates a deck with business logic validation
func (s *DeckService) Update(id uuid.UUID, req *models.UpdateDeckRequest) (*models.Deck, error) {
	// Get existing deck first
	existingDeck, err := s.deckRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("deck not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) == 0 {
		return existingDeck, nil // No changes needed
	}

	updatedDeck, err := s.deckRepo.Update(id, updates)
	if err != nil {
		s.Logger.WithError(err).WithField("deck_id", id).Error("Service failed to update deck")
		return nil, fmt.Errorf("failed to update deck: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"deck_id": id,
	}).Info("Deck updated successfully")

	return updatedDeck, nil
}

// UpdateWithOwnership updates a deck with user ownership validation
func (s *DeckService) UpdateWithOwnership(id uuid.UUID, userID uuid.UUID, req *models.UpdateDeckRequest) (*models.Deck, error) {
	// Get existing deck first
	existingDeck, err := s.deckRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("deck not found: %w", err)
	}

	// Check if deck belongs to user
	if existingDeck.UserID != userID {
		s.Logger.WithFields(logrus.Fields{
			"deck_id":  id,
			"user_id":  userID,
			"owner_id": existingDeck.UserID,
		}).Warn("Unauthorized attempt to update deck")
		return nil, fmt.Errorf("unauthorized: deck does not belong to user")
	}

	// Call the regular update method
	return s.Update(id, req)
}

// Delete removes a deck with validation
func (s *DeckService) Delete(id uuid.UUID) error {
	// Check if deck exists first
	_, err := s.deckRepo.GetByID(id)
	if err != nil {
		s.Logger.WithField("deck_id", id).Warn("Attempted to delete non-existent deck")
		return fmt.Errorf("deck not found: %w", err)
	}

	// Check if deck has flashcards (optional business logic)
	flashcardCount, err := s.deckRepo.GetDeckFlashcardCount(id)
	if err != nil {
		s.Logger.WithError(err).WithField("deck_id", id).Error("Failed to check deck flashcard count")
		return fmt.Errorf("failed to check deck contents: %w", err)
	}

	if flashcardCount > 0 {
		s.Logger.WithFields(logrus.Fields{
			"deck_id":         id,
			"flashcard_count": flashcardCount,
		}).Warn("Deleting deck with flashcards")
	}

	err = s.deckRepo.Delete(id)
	if err != nil {
		s.Logger.WithError(err).WithField("deck_id", id).Error("Service failed to delete deck")
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"deck_id":         id,
		"flashcard_count": flashcardCount,
	}).Info("Deck deleted successfully")

	return nil
}

// DeleteWithOwnership removes a deck with user ownership validation
func (s *DeckService) DeleteWithOwnership(id uuid.UUID, userID uuid.UUID) error {
	// Get existing deck first
	existingDeck, err := s.deckRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("deck not found: %w", err)
	}

	// Check if deck belongs to user
	if existingDeck.UserID != userID {
		s.Logger.WithFields(logrus.Fields{
			"deck_id":  id,
			"user_id":  userID,
			"owner_id": existingDeck.UserID,
		}).Warn("Unauthorized attempt to delete deck")
		return fmt.Errorf("unauthorized: deck does not belong to user")
	}

	// Call regular delete method
	return s.Delete(id)
}
