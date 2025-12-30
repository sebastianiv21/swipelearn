package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"swipelearn-api/internal/models"
	"swipelearn-api/pkg/testutils"
)

// MockFlashcardRepository is a mock implementation of FlashcardRepository for testing
type MockFlashcardRepository struct {
	mock.Mock
}

func (m *MockFlashcardRepository) Create(card *models.Flashcard) (*models.Flashcard, error) {
	args := m.Called(card)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flashcard), args.Error(1)
}

func (m *MockFlashcardRepository) GetByID(id uuid.UUID) (*models.Flashcard, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flashcard), args.Error(1)
}

func (m *MockFlashcardRepository) GetByUser(userID uuid.UUID) ([]*models.Flashcard, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Flashcard), args.Error(1)
}

func (m *MockFlashcardRepository) Update(id uuid.UUID, updates *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flashcard), args.Error(1)
}

func (m *MockFlashcardRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestFlashcardService_Create_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	userID := uuid.New()
	deckID := uuid.New()
	req := &models.CreateFlashcardRequest{
		Front:  "Question",
		Back:   "Answer",
		UserID: userID,
		DeckID: deckID,
	}

	expectedCard := &models.Flashcard{
		ID:          uuid.New(),
		UserID:      userID,
		DeckID:      deckID,
		Front:       req.Front,
		Back:        req.Back,
		Difficulty:  2.5,
		Interval:    1,
		EaseFactor:  2.5,
		ReviewCount: 0,
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Flashcard")).Return(expectedCard, nil)

	result, err := service.Create(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedCard.ID, result.ID)
	assert.Equal(t, expectedCard.UserID, result.UserID)
	assert.Equal(t, expectedCard.DeckID, result.DeckID)
	assert.Equal(t, expectedCard.Front, result.Front)
	assert.Equal(t, expectedCard.Back, result.Back)
	assert.Equal(t, 2.5, result.Difficulty)
	assert.Equal(t, 1, result.Interval)
	assert.Equal(t, 2.5, result.EaseFactor)
	assert.Equal(t, 0, result.ReviewCount)

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_Create_InvalidUserID(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	deckID := uuid.New()
	req := &models.CreateFlashcardRequest{
		Front:  "Question",
		Back:   "Answer",
		UserID: uuid.Nil, // Invalid
		DeckID: deckID,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user ID is required")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestFlashcardService_Create_InvalidDeckID(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	userID := uuid.New()
	req := &models.CreateFlashcardRequest{
		Front:  "Question",
		Back:   "Answer",
		UserID: userID,
		DeckID: uuid.Nil, // Invalid
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deck ID is required")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestFlashcardService_GetByUser_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	userID := uuid.New()
	expectedCards := []*models.Flashcard{
		{
			ID:          uuid.New(),
			UserID:      userID,
			DeckID:      uuid.New(),
			Front:       "Question 1",
			Back:        "Answer 1",
			Difficulty:  2.5,
			Interval:    1,
			EaseFactor:  2.5,
			ReviewCount: 0,
		},
	}

	mockRepo.On("GetByUser", userID).Return(expectedCards, nil)

	result, err := service.GetByUser(userID, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, expectedCards[0].ID, result[0].ID)

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_ReviewFlashcard_PerfectResponse(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	quality := 5 // Perfect response

	now := time.Now()
	existingCard := &models.Flashcard{
		ID:          cardID,
		UserID:      uuid.New(),
		DeckID:      uuid.New(),
		Front:       "Question",
		Back:        "Answer",
		Difficulty:  2.5,
		Interval:    1,
		EaseFactor:  2.5,
		ReviewCount: 0,
		LastReview:  nil,
		NextReview:  nil,
	}

	// Expected card after perfect response (quality = 5)
	expectedCard := &models.Flashcard{
		ID:          cardID,
		UserID:      existingCard.UserID,
		DeckID:      existingCard.DeckID,
		Front:       existingCard.Front,
		Back:        existingCard.Back,
		Difficulty:  2.6, // EF' = 2.5 + (0.1 - (5-5)*(0.08+(5-5)*0.02)) = 2.5 + 0.1 = 2.6
		Interval:    6,   // First correct review: 6 days
		EaseFactor:  2.6, // Should match difficulty for SM-2
		ReviewCount: 1,
		LastReview:  &now,
		NextReview:  &time.Time{}, // Will be set in test
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)
	// Mock Update returns updated card
	mockRepo.On("Update", cardID, mock.AnythingOfType("*models.UpdateFlashcardRequest")).Return(expectedCard, nil)

	result, err := service.ReviewFlashcard(cardID, quality)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedCard.ID, result.ID)
	assert.Equal(t, 2.6, result.Difficulty) // New ease factor
	assert.Equal(t, 6, result.Interval)     // New interval
	assert.Equal(t, 2.6, result.EaseFactor) // Should match difficulty
	assert.Equal(t, 1, result.ReviewCount)  // Incremented
	assert.NotNil(t, result.LastReview)
	assert.NotNil(t, result.NextReview)

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_ReviewFlashcard_PoorResponse(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	quality := 2 // Poor response (below threshold)

	now := time.Now()
	existingCard := &models.Flashcard{
		ID:          cardID,
		UserID:      uuid.New(),
		DeckID:      uuid.New(),
		Front:       "Question",
		Back:        "Answer",
		Difficulty:  2.5,
		Interval:    6,
		EaseFactor:  2.5,
		ReviewCount: 2,
	}

	// Expected card after poor response (quality < 3)
	expectedCard := &models.Flashcard{
		ID:          cardID,
		UserID:      existingCard.UserID,
		DeckID:      existingCard.DeckID,
		Front:       existingCard.Front,
		Back:        existingCard.Back,
		Difficulty:  1.3, // Minimum ease factor
		Interval:    1,   // Reset to 1 day
		EaseFactor:  1.3, // Minimum ease factor
		ReviewCount: 3,   // Incremented
		LastReview:  &now,
		NextReview:  &time.Time{}, // Will be set in test
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)
	// Mock Update returns updated card
	mockRepo.On("Update", cardID, mock.AnythingOfType("*models.UpdateFlashcardRequest")).Return(expectedCard, nil)

	result, err := service.ReviewFlashcard(cardID, quality)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1.3, result.Difficulty) // Minimum ease factor
	assert.Equal(t, 1, result.Interval)     // Reset interval
	assert.Equal(t, 1.3, result.EaseFactor) // Minimum ease factor
	assert.Equal(t, 3, result.ReviewCount)  // Incremented

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_ReviewFlashcard_InvalidQuality(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	quality := 6 // Invalid (must be 0-5)

	result, err := service.ReviewFlashcard(cardID, quality)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quality must be between 0 and 5")

	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestFlashcardService_ReviewFlashcard_CardNotFound(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	quality := 3

	mockRepo.On("GetByID", cardID).Return(nil, sql.ErrNoRows)

	result, err := service.ReviewFlashcard(cardID, quality)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "flashcard not found")

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_GetDueCards_Empty(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	userID := uuid.New()

	// Mock cards that are not due
	tomorrow := time.Now().Add(24 * time.Hour)
	cardsNotDue := []*models.Flashcard{
		{
			ID:         uuid.New(),
			UserID:     userID,
			NextReview: &tomorrow, // Due tomorrow
		},
	}

	mockRepo.On("GetByUser", userID).Return(cardsNotDue, nil)

	result, err := service.GetDueCards(userID)

	require.NoError(t, err)
	assert.Empty(t, result) // No cards due

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_GetDueCards_WithDueCards(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	userID := uuid.New()
	now := time.Now()

	// Mock cards with different due dates
	oneHourAgo := now.Add(-1 * time.Hour)
	oneHourFromNow := now.Add(1 * time.Hour)

	cards := []*models.Flashcard{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Front:      "Due Card 1",
			NextReview: &oneHourAgo, // Due 1 hour ago
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Front:      "Due Card 2",
			NextReview: nil, // Never reviewed, so it's due
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Front:      "Future Card",
			NextReview: &oneHourFromNow, // Due in 1 hour
		},
	}

	mockRepo.On("GetByUser", userID).Return(cards, nil)

	result, err := service.GetDueCards(userID)

	require.NoError(t, err)
	assert.Len(t, result, 2) // Only 2 cards are due

	// Verify the due cards are returned
	fronts := []string{result[0].Front, result[1].Front}
	assert.Contains(t, fronts, "Due Card 1")
	assert.Contains(t, fronts, "Due Card 2")
	assert.NotContains(t, fronts, "Future Card")

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_ReviewFlashcardWithOwnership_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	userID := uuid.New()
	quality := 4 // Good response

	existingCard := &models.Flashcard{
		ID:          cardID,
		UserID:      userID,
		Front:       "Question",
		Back:        "Answer",
		Difficulty:  2.5,
		Interval:    1,
		EaseFactor:  2.5,
		ReviewCount: 0,
	}

	expectedCard := &models.Flashcard{
		ID:          cardID,
		UserID:      userID,
		Front:       "Question",
		Back:        "Answer",
		Difficulty:  2.5,
		Interval:    6, // Second correct review
		EaseFactor:  2.5,
		ReviewCount: 1,
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)
	// Mock Update returns updated card
	mockRepo.On("Update", cardID, mock.AnythingOfType("*models.UpdateFlashcardRequest")).Return(expectedCard, nil)

	result, err := service.ReviewFlashcardWithOwnership(cardID, userID, quality)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedCard.ID, result.ID)
	assert.Equal(t, expectedCard.ReviewCount, result.ReviewCount)

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_ReviewFlashcardWithOwnership_Unauthorized(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()
	quality := 4

	existingCard := &models.Flashcard{
		ID:     cardID,
		UserID: ownerID, // Different user
		Front:  "Question",
		Back:   "Answer",
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)

	result, err := service.ReviewFlashcardWithOwnership(cardID, userID, quality)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Contains(t, err.Error(), "does not belong to user")

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_UpdateWithOwnership_Unauthorized(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	existingCard := &models.Flashcard{
		ID:     cardID,
		UserID: ownerID, // Different user
		Front:  "Question",
		Back:   "Answer",
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)

	req := &models.UpdateFlashcardRequest{
		Front: func() *string { s := "New Question"; return &s }(),
	}

	result, err := service.UpdateWithOwnership(cardID, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Contains(t, err.Error(), "does not belong to user")

	mockRepo.AssertExpectations(t)
}

func TestFlashcardService_DeleteWithOwnership_Unauthorized(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockFlashcardRepository{}
	service := NewFlashcardService(mockRepo, logger)

	cardID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	existingCard := &models.Flashcard{
		ID:     cardID,
		UserID: ownerID, // Different user
		Front:  "Question",
		Back:   "Answer",
	}

	// Mock GetByID returns existing card
	mockRepo.On("GetByID", cardID).Return(existingCard, nil)

	err := service.DeleteWithOwnership(cardID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Contains(t, err.Error(), "does not belong to user")

	mockRepo.AssertExpectations(t)
}
