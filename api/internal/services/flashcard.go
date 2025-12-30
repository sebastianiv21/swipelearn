package services

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"swipelearn-api/internal/models"
	"swipelearn-api/internal/repositories"
)

type FlashcardService struct {
	flashcardRepo repositories.FlashcardRepositoryInterface
	Logger        *logrus.Logger
}

func NewFlashcardService(repo repositories.FlashcardRepositoryInterface, logger *logrus.Logger) *FlashcardService {
	return &FlashcardService{
		flashcardRepo: repo,
		Logger:        logger,
	}
}

// Create creates a new flashcard with business logic validation
func (s *FlashcardService) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	// Business logic validation
	if req.DeckID == uuid.Nil {
		return nil, fmt.Errorf("deck ID is required")
	}
	if req.UserID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	card := &models.Flashcard{
		ID:          uuid.New(),
		UserID:      req.UserID,
		Front:       req.Front,
		Back:        req.Back,
		DeckID:      req.DeckID,
		Difficulty:  2.5, // Initial difficulty for new cards
		Interval:    1,   // Start with 1 day interval
		EaseFactor:  2.5, // SM-2 default ease factor
		ReviewCount: 0,
	}

	savedCard, err := s.flashcardRepo.Create(card)
	if err != nil {
		s.Logger.WithError(err).Error("Service failed to create flashcard")
		return nil, fmt.Errorf("failed to create flashcard: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"flashcard_id": savedCard.ID,
		"user_id":      savedCard.UserID,
		"deck_id":      savedCard.DeckID,
	}).Info("Flashcard created successfully")

	return savedCard, nil
}

// GetByUser retrieves flashcards for a user with optional filters
func (s *FlashcardService) GetByUser(userID uuid.UUID, filters map[string]any) ([]*models.Flashcard, error) {
	flashcards, err := s.flashcardRepo.GetByUser(userID)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", userID).Error("Service failed to get flashcards")
		return nil, fmt.Errorf("failed to get flashcards: %w", err)
	}

	// Apply business logic filters if provided
	if filters != nil {
		// Example: filter by difficulty
		if minDifficulty, ok := filters["min_difficulty"].(float64); ok {
			var filtered []*models.Flashcard
			for _, card := range flashcards {
				if card.Difficulty >= minDifficulty {
					filtered = append(filtered, card)
				}
			}
			flashcards = filtered
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id":         userID,
		"flashcard_count": len(flashcards),
		"filters":         filters,
	}).Info("Retrieved flashcards for user")

	return flashcards, nil
}

// Update updates a flashcard with spaced repetition logic
func (s *FlashcardService) Update(id uuid.UUID, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	// Get existing card first
	existingCard, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// Apply spaced repetition algorithm updates
	if req.Difficulty != nil {
		// SM-2 algorithm: new ease factor = EF + (0.1 - (5 - q) * (EF + q))
		// where q = response quality (0-5), EF = ease factor
		// Simplified: if correct (q=5), increase EF slightly
		q := 3.0 // Assume average response quality
		*req.Difficulty = existingCard.EaseFactor + (0.1 - (5-q)*(existingCard.EaseFactor+q))

		// Adjust interval based on new difficulty
		*req.Difficulty = math.Max(1.3, *req.Difficulty) // Minimum ease factor
	}

	updatedCard, err := s.flashcardRepo.Update(id, req)
	if err != nil {
		s.Logger.WithError(err).WithField("flashcard_id", id).Error("Service failed to update flashcard")
		return nil, fmt.Errorf("failed to update flashcard: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"flashcard_id":   id,
		"new_difficulty": req.Difficulty,
	}).Info("Flashcard updated successfully")

	return updatedCard, nil
}

// UpdateWithOwnership updates a flashcard with user ownership validation
func (s *FlashcardService) UpdateWithOwnership(id uuid.UUID, userID uuid.UUID, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	// Get existing card first
	existingCard, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// Check if the card belongs to the user
	if existingCard.UserID != userID {
		return nil, fmt.Errorf("unauthorized: flashcard does not belong to user")
	}

	// Call the regular update method
	return s.Update(id, req)
}

// Delete removes a flashcard with validation
func (s *FlashcardService) Delete(id uuid.UUID) error {
	// Check if card exists first
	_, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		s.Logger.WithField("flashcard_id", id).Warn("Attempted to delete non-existent flashcard")
		return fmt.Errorf("flashcard not found: %w", err)
	}

	err = s.flashcardRepo.Delete(id)
	if err != nil {
		s.Logger.WithError(err).WithField("flashcard_id", id).Error("Service failed to delete flashcard")
		return fmt.Errorf("failed to delete flashcard: %w", err)
	}

	s.Logger.WithField("flashcard_id", id).Info("Flashcard deleted successfully")
	return nil
}

// DeleteWithOwnership removes a flashcard with user ownership validation
func (s *FlashcardService) DeleteWithOwnership(id uuid.UUID, userID uuid.UUID) error {
	// Get existing card first
	existingCard, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("flashcard not found: %w", err)
	}

	// Check if the card belongs to the user
	if existingCard.UserID != userID {
		return fmt.Errorf("unauthorized: flashcard does not belong to user")
	}

	// Call the regular delete method
	return s.Delete(id)
}

// ReviewFlashcard handles the spaced repetition review logic using correct SM-2 algorithm
func (s *FlashcardService) ReviewFlashcard(id uuid.UUID, quality int) (*models.Flashcard, error) {
	// Validate quality range (0-5)
	if quality < 0 || quality > 5 {
		return nil, fmt.Errorf("quality must be between 0 and 5, got %d", quality)
	}

	card, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// SM-2 Algorithm - Correct Formula
	q := float64(quality)

	// Correct SM-2 ease factor formula:
	// EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
	newEaseFactor := card.EaseFactor + (0.1 - (5.0-q)*(0.08+(5.0-q)*0.02))

	// Enforce minimum ease factor of 1.3
	newEaseFactor = math.Max(1.3, newEaseFactor)

	var newInterval int
	var newRepetitions int
	var nextReview time.Time

	if q < 3 {
		// Incorrect response (quality 0, 1, or 2), reset interval and repetitions
		newInterval = 1
		newRepetitions = 0
		nextReview = time.Now().Add(time.Hour * 24)
	} else {
		// Correct response (quality 3, 4, or 5)
		newRepetitions = card.ReviewCount + 1

		// Calculate new interval based on repetitions
		switch newRepetitions {
		case 1:
			newInterval = 1
		case 2:
			newInterval = 6
		default:
			newInterval = int(math.Round(float64(card.Interval) * newEaseFactor))
		}
		nextReview = time.Now().Add(time.Hour * 24 * time.Duration(newInterval))
	}

	// Update the card with all SM-2 fields
	updateReq := &models.UpdateFlashcardRequest{
		Difficulty:  &newEaseFactor,
		Interval:    &newInterval,
		EaseFactor:  &newEaseFactor,
		ReviewCount: &newRepetitions,
		LastReview:  &[]time.Time{time.Now()}[0],
		NextReview:  &nextReview,
	}

	updatedCard, err := s.flashcardRepo.Update(id, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update flashcard review: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"flashcard_id":    id,
		"quality":         quality,
		"new_interval":    newInterval,
		"new_ease_factor": newEaseFactor,
		"repetitions":     newRepetitions,
		"next_review":     nextReview,
	}).Info("Flashcard reviewed successfully with SM-2 algorithm")

	return updatedCard, nil
}

// ReviewFlashcardWithOwnership handles the spaced repetition review logic with user ownership validation
func (s *FlashcardService) ReviewFlashcardWithOwnership(id uuid.UUID, userID uuid.UUID, quality int) (*models.Flashcard, error) {
	// Get the flashcard first
	card, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// Check if the card belongs to the user
	if card.UserID != userID {
		s.Logger.WithFields(logrus.Fields{
			"flashcard_id": id,
			"user_id":      userID,
			"owner_id":     card.UserID,
		}).Warn("Unauthorized attempt to review flashcard")
		return nil, fmt.Errorf("unauthorized: flashcard does not belong to user")
	}

	// Call the regular review method
	return s.ReviewFlashcard(id, quality)
}

// GetDueCards retrieves flashcards that are due for review
func (s *FlashcardService) GetDueCards(userID uuid.UUID) ([]*models.Flashcard, error) {
	flashcards, err := s.flashcardRepo.GetByUser(userID)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", userID).Error("Service failed to get flashcards for due cards")
		return nil, fmt.Errorf("failed to get flashcards: %w", err)
	}

	var dueCards []*models.Flashcard
	now := time.Now()

	for _, card := range flashcards {
		// If next_review is nil or is in the past, card is due
		if card.NextReview == nil || card.NextReview.Before(now) {
			dueCards = append(dueCards, card)
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id":        userID,
		"due_card_count": len(dueCards),
	}).Info("Retrieved due flashcards for user")

	return dueCards, nil
}
