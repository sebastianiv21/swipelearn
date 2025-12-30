package services

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"swipelearn-api/internal/models"
	"swipelearn-api/pkg/testutils"
)

// MockDeckRepository is a mock implementation of DeckRepository for testing
type MockDeckRepository struct {
	mock.Mock
}

func (m *MockDeckRepository) Create(deck *models.Deck) (*models.Deck, error) {
	args := m.Called(deck)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Deck), args.Error(1)
}

func (m *MockDeckRepository) GetByID(id uuid.UUID) (*models.Deck, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Deck), args.Error(1)
}

func (m *MockDeckRepository) GetAll() ([]*models.Deck, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Deck), args.Error(1)
}

func (m *MockDeckRepository) GetByUser(userID uuid.UUID) ([]*models.Deck, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Deck), args.Error(1)
}

func (m *MockDeckRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.Deck, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Deck), args.Error(1)
}

func (m *MockDeckRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDeckRepository) GetDeckFlashcardCount(deckID uuid.UUID) (int, error) {
	args := m.Called(deckID)
	return args.Int(0), args.Error(1)
}

func TestDeckService_Create_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	userID := uuid.New()
	req := &models.CreateDeckRequest{
		Name:        "Test Deck",
		Description: "Test Description",
	}

	expectedDeck := &models.Deck{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Deck")).Return(expectedDeck, nil)

	result, err := service.Create(req, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedDeck.ID, result.ID)
	assert.Equal(t, expectedDeck.UserID, result.UserID)
	assert.Equal(t, expectedDeck.Name, result.Name)
	assert.Equal(t, expectedDeck.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Create_RepositoryError(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	userID := uuid.New()
	req := &models.CreateDeckRequest{
		Name:        "Test Deck",
		Description: "Test Description",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Deck")).Return(nil, assert.AnError)

	result, err := service.Create(req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create deck")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetByID_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	expectedDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Test Deck",
		Description: "Test Description",
	}

	mockRepo.On("GetByID", deckID).Return(expectedDeck, nil)

	result, err := service.GetByID(deckID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedDeck.ID, result.ID)
	assert.Equal(t, expectedDeck.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetByID_NotFound(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	mockRepo.On("GetByID", deckID).Return(nil, sql.ErrNoRows)

	result, err := service.GetByID(deckID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get deck")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetByIDWithOwnership_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	userID := uuid.New()
	expectedDeck := &models.Deck{
		ID:          deckID,
		UserID:      userID,
		Name:        "Test Deck",
		Description: "Test Description",
	}

	mockRepo.On("GetByID", deckID).Return(expectedDeck, nil)

	result, err := service.GetByIDWithOwnership(deckID, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedDeck.ID, result.ID)
	assert.Equal(t, expectedDeck.UserID, result.UserID)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetByIDWithOwnership_Unauthorized(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	deck := &models.Deck{
		ID:          deckID,
		UserID:      ownerID, // Different user
		Name:        "Test Deck",
		Description: "Test Description",
	}

	mockRepo.On("GetByID", deckID).Return(deck, nil)

	result, err := service.GetByIDWithOwnership(deckID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Contains(t, err.Error(), "does not belong to user")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetAll_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	expectedDecks := []*models.Deck{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "Deck 1",
			Description: "Description 1",
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "Deck 2",
			Description: "Description 2",
		},
	}

	mockRepo.On("GetAll").Return(expectedDecks, nil)

	result, err := service.GetAll()

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, expectedDecks[0].ID, result[0].ID)
	assert.Equal(t, expectedDecks[1].ID, result[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetAll_RepositoryError(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	mockRepo.On("GetAll").Return(nil, assert.AnError)

	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get decks")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_GetByUser_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	userID := uuid.New()
	expectedDecks := []*models.Deck{
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Deck 1",
			Description: "Description 1",
		},
	}

	mockRepo.On("GetByUser", userID).Return(expectedDecks, nil)

	result, err := service.GetByUser(userID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, expectedDecks[0].ID, result[0].ID)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Update_Name(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	newName := "Updated Deck Name"

	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Original Name",
		Description: "Original Description",
	}

	updatedDeck := &models.Deck{
		ID:          deckID,
		UserID:      existingDeck.UserID,
		Name:        newName,
		Description: "Original Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock Update returns updated deck
	mockRepo.On("Update", deckID, map[string]interface{}{"name": newName}).Return(updatedDeck, nil)

	req := &models.UpdateDeckRequest{
		Name: &newName,
	}

	result, err := service.Update(deckID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedDeck.ID, result.ID)
	assert.Equal(t, updatedDeck.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Update_Description(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	newDescription := "Updated Description"

	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Original Name",
		Description: "Original Description",
	}

	updatedDeck := &models.Deck{
		ID:          deckID,
		UserID:      existingDeck.UserID,
		Name:        "Original Name",
		Description: newDescription,
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock Update returns updated deck
	mockRepo.On("Update", deckID, map[string]interface{}{"description": newDescription}).Return(updatedDeck, nil)

	req := &models.UpdateDeckRequest{
		Description: &newDescription,
	}

	result, err := service.Update(deckID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedDeck.ID, result.ID)
	assert.Equal(t, updatedDeck.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Update_NoChanges(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	originalName := "Original Name"

	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        originalName,
		Description: "Original Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)

	req := &models.UpdateDeckRequest{
		Name: &originalName, // Same name as existing
	}

	result, err := service.Update(deckID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, existingDeck.ID, result.ID)
	assert.Equal(t, existingDeck.Name, result.Name)

	// No Update call should be made since there are no changes
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestDeckService_UpdateWithOwnership_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	userID := uuid.New()
	newName := "Updated Deck Name"

	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      userID,
		Name:        "Original Name",
		Description: "Original Description",
	}

	updatedDeck := &models.Deck{
		ID:          deckID,
		UserID:      userID,
		Name:        newName,
		Description: "Original Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock Update returns updated deck
	mockRepo.On("Update", deckID, map[string]interface{}{"name": newName}).Return(updatedDeck, nil)

	req := &models.UpdateDeckRequest{
		Name: &newName,
	}

	result, err := service.UpdateWithOwnership(deckID, userID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedDeck.ID, result.ID)
	assert.Equal(t, updatedDeck.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestDeckService_UpdateWithOwnership_Unauthorized(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      ownerID, // Different user
		Name:        "Original Name",
		Description: "Original Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)

	req := &models.UpdateDeckRequest{
		Name: func() *string { s := "Updated Name"; return &s }(),
	}

	result, err := service.UpdateWithOwnership(deckID, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Contains(t, err.Error(), "does not belong to user")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Delete_Success(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Test Deck",
		Description: "Test Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock GetDeckFlashcardCount returns 0 (no flashcards)
	mockRepo.On("GetDeckFlashcardCount", deckID).Return(0, nil)
	// Mock Delete returns success
	mockRepo.On("Delete", deckID).Return(nil)

	err := service.Delete(deckID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeckService_Delete_WithFlashcards(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Test Deck",
		Description: "Test Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock GetDeckFlashcardCount returns 5 (has flashcards)
	mockRepo.On("GetDeckFlashcardCount", deckID).Return(5, nil)
	// Mock Delete returns success
	mockRepo.On("Delete", deckID).Return(nil)

	err := service.Delete(deckID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeckService_Delete_NotFound(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()

	// Mock GetByID returns not found
	mockRepo.On("GetByID", deckID).Return(nil, sql.ErrNoRows)

	err := service.Delete(deckID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deck not found")

	mockRepo.AssertExpectations(t)
}

func TestDeckService_Delete_RepositoryError(t *testing.T) {
	logger := testutils.TestLogger()
	mockRepo := &MockDeckRepository{}
	service := NewDeckService(mockRepo, logger)

	deckID := uuid.New()
	existingDeck := &models.Deck{
		ID:          deckID,
		UserID:      uuid.New(),
		Name:        "Test Deck",
		Description: "Test Description",
	}

	// Mock GetByID returns existing deck
	mockRepo.On("GetByID", deckID).Return(existingDeck, nil)
	// Mock GetDeckFlashcardCount returns 0
	mockRepo.On("GetDeckFlashcardCount", deckID).Return(0, nil)
	// Mock Delete returns error
	mockRepo.On("Delete", deckID).Return(assert.AnError)

	err := service.Delete(deckID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete deck")

	mockRepo.AssertExpectations(t)
}
