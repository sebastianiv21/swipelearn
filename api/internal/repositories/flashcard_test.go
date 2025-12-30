package repositories

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swipelearn-api/internal/models"
	"swipelearn-api/pkg/testutils"
)

func TestFlashcardRepository_Create_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user and deck first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)

	flashcard := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)

	createdFlashcard, err := repo.Create(flashcard)
	require.NoError(t, err)
	require.NotNil(t, createdFlashcard)

	assert.NotEqual(t, uuid.Nil, createdFlashcard.ID)
	assert.Equal(t, flashcard.UserID, createdFlashcard.UserID)
	assert.Equal(t, flashcard.DeckID, createdFlashcard.DeckID)
	assert.Equal(t, flashcard.Front, createdFlashcard.Front)
	assert.Equal(t, flashcard.Back, createdFlashcard.Back)
	assert.Equal(t, 2.5, createdFlashcard.Difficulty)
	assert.Equal(t, 1, createdFlashcard.Interval)
	assert.Equal(t, 2.5, createdFlashcard.EaseFactor)
	assert.Equal(t, 0, createdFlashcard.ReviewCount)
	assert.False(t, createdFlashcard.CreatedAt.IsZero())
	assert.False(t, createdFlashcard.UpdatedAt.IsZero())
}

func TestFlashcardRepository_GetByID_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user, deck, and flashcard
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)
	flashcard := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	createdFlashcard, err := repo.Create(flashcard)
	require.NoError(t, err)

	// Get flashcard by ID
	retrievedFlashcard, err := repo.GetByID(createdFlashcard.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedFlashcard)

	assert.Equal(t, createdFlashcard.ID, retrievedFlashcard.ID)
	assert.Equal(t, createdFlashcard.UserID, retrievedFlashcard.UserID)
	assert.Equal(t, createdFlashcard.DeckID, retrievedFlashcard.DeckID)
	assert.Equal(t, createdFlashcard.Front, retrievedFlashcard.Front)
	assert.Equal(t, createdFlashcard.Back, retrievedFlashcard.Back)
}

func TestFlashcardRepository_GetByID_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)

	// Try to get a non-existent flashcard
	randomID := uuid.New()
	flashcard, err := repo.GetByID(randomID)

	assert.Error(t, err)
	assert.Nil(t, flashcard)
	assert.Contains(t, err.Error(), "flashcard not found")
}

func TestFlashcardRepository_GetByUser_Empty(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)

	// Get all flashcards for a user when table is empty
	flashcards, err := repo.GetByUser(uuid.New())
	require.NoError(t, err)
	assert.Empty(t, flashcards)
}

func TestFlashcardRepository_GetByUser_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user and deck
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)

	// Create multiple flashcards
	flashcard1 := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	flashcard1.Front = "Question 1"
	flashcard2 := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	flashcard2.Front = "Question 2"

	_, err = repo.Create(flashcard1)
	require.NoError(t, err)
	_, err = repo.Create(flashcard2)
	require.NoError(t, err)

	// Get flashcards for user
	flashcards, err := repo.GetByUser(createdUser.ID)
	require.NoError(t, err)
	assert.Len(t, flashcards, 2)

	// Verify flashcards are ordered by created_at DESC (newest first)
	if flashcards[0].Front == "Question 1" {
		assert.Equal(t, "Question 1", flashcards[0].Front)
		assert.Equal(t, "Question 2", flashcards[1].Front)
	} else {
		assert.Equal(t, "Question 1", flashcards[1].Front)
		assert.Equal(t, "Question 2", flashcards[0].Front)
	}
}

func TestFlashcardRepository_Update_AllFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user, deck, and flashcard
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)
	flashcard := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	createdFlashcard, err := repo.Create(flashcard)
	require.NoError(t, err)

	// Update all fields
	newFront := "Updated Question"
	newBack := "Updated Answer"
	newDifficulty := 2.8
	newInterval := 3
	newEaseFactor := 2.7
	newReviewCount := 5
	newLastReview := time.Now()
	newNextReview := time.Now().Add(3 * 24 * time.Hour)

	updates := &models.UpdateFlashcardRequest{
		Front:       &newFront,
		Back:        &newBack,
		Difficulty:  &newDifficulty,
		Interval:    &newInterval,
		EaseFactor:  &newEaseFactor,
		ReviewCount: &newReviewCount,
		LastReview:  &newLastReview,
		NextReview:  &newNextReview,
	}

	updatedFlashcard, err := repo.Update(createdFlashcard.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedFlashcard)

	assert.Equal(t, createdFlashcard.ID, updatedFlashcard.ID)
	assert.Equal(t, newFront, updatedFlashcard.Front)
	assert.Equal(t, newBack, updatedFlashcard.Back)
	assert.Equal(t, newDifficulty, updatedFlashcard.Difficulty)
	assert.Equal(t, newInterval, updatedFlashcard.Interval)
	assert.Equal(t, newEaseFactor, updatedFlashcard.EaseFactor)
	assert.Equal(t, newReviewCount, updatedFlashcard.ReviewCount)
	assert.NotNil(t, updatedFlashcard.LastReview)
	assert.NotNil(t, updatedFlashcard.NextReview)
	assert.True(t, updatedFlashcard.UpdatedAt.After(createdFlashcard.UpdatedAt))
}

func TestFlashcardRepository_Update_PartialFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user, deck, and flashcard
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)
	flashcard := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	createdFlashcard, err := repo.Create(flashcard)
	require.NoError(t, err)

	// Update only front field
	newFront := "Updated Question Only"
	updates := &models.UpdateFlashcardRequest{
		Front: &newFront,
	}

	updatedFlashcard, err := repo.Update(createdFlashcard.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedFlashcard)

	assert.Equal(t, createdFlashcard.ID, updatedFlashcard.ID)
	assert.Equal(t, newFront, updatedFlashcard.Front)
	// Other fields should remain unchanged
	assert.Equal(t, createdFlashcard.Back, updatedFlashcard.Back)
	assert.Equal(t, createdFlashcard.Difficulty, updatedFlashcard.Difficulty)
}

func TestFlashcardRepository_Delete_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user, deck, and flashcard
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	deckRepo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := deckRepo.Create(deck)
	require.NoError(t, err)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)
	flashcard := testutils.CreateTestFlashcard(createdUser.ID, createdDeck.ID)
	createdFlashcard, err := repo.Create(flashcard)
	require.NoError(t, err)

	// Delete flashcard
	err = repo.Delete(createdFlashcard.ID)
	require.NoError(t, err)

	// Verify flashcard is deleted
	deletedFlashcard, err := repo.GetByID(createdFlashcard.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedFlashcard)
	assert.Contains(t, err.Error(), "flashcard not found")
}

func TestFlashcardRepository_Delete_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewFlashcardRepository(td.DB.DB, td.Logger)

	// Try to delete a non-existent flashcard
	randomID := uuid.New()
	err := repo.Delete(randomID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "flashcard not found")
}
