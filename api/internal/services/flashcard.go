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
	flashcardRepo *repositories.FlashcardRepository
	Logger        *logrus.Logger
}

func NewFlashcardService(repo *repositories.FlashcardRepository, logger *logrus.Logger) *FlashcardService {
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

	card := &models.Flashcard{
		ID:          uuid.New(),
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

// ReviewFlashcard handles the spaced repetition review logic
func (s *FlashcardService) ReviewFlashcard(id uuid.UUID, quality int) (*models.Flashcard, error) {
	card, err := s.flashcardRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("flashcard not found: %w", err)
	}

	// SM-2 Algorithm
	q := float64(quality)
	easeFactor := card.EaseFactor + (0.1 - (5-q)*(card.EaseFactor+q))

	var newInterval int
	var nextReview time.Time

	if q < 3 {
		// Incorrect response, reset interval
		newInterval = 1
		nextReview = time.Now().Add(time.Hour * 24)
	} else {
		// Correct response, calculate new interval
		switch card.ReviewCount {
		case 0:
			newInterval = 1
		case 1:
			newInterval = 6
		default:
			newInterval = int(math.Round(float64(card.Interval) * easeFactor))
		}
		nextReview = time.Now().Add(time.Hour * 24 * time.Duration(newInterval))
	}

	// Update the card
	updateReq := &models.UpdateFlashcardRequest{
		Difficulty: &easeFactor,
		Interval:   &newInterval,
	}

	// Store review record
	card.ReviewCount++
	card.LastReview = &time.Time{}
	*card.LastReview = time.Now()
	card.NextReview = &nextReview

	updatedCard, err := s.flashcardRepo.Update(id, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update flashcard review: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"flashcard_id": id,
		"quality":      quality,
		"new_interval": newInterval,
		"next_review":  nextReview,
	}).Info("Flashcard reviewed successfully")

	return updatedCard, nil
}
