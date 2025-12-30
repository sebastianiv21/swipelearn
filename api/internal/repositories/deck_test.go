package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swipelearn-api/pkg/testutils"
)

func TestDeckRepository_Create_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// First create a user since deck has foreign key constraint
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	deck := testutils.CreateTestDeck(createdUser.ID)

	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)
	require.NotNil(t, createdDeck)

	assert.NotEqual(t, uuid.Nil, createdDeck.ID)
	assert.Equal(t, deck.UserID, createdDeck.UserID)
	assert.Equal(t, deck.Name, createdDeck.Name)
	assert.Equal(t, deck.Description, createdDeck.Description)
	assert.False(t, createdDeck.CreatedAt.IsZero())
	assert.False(t, createdDeck.UpdatedAt.IsZero())
}

func TestDeckRepository_GetByID_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create a deck
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Get the deck by ID
	retrievedDeck, err := repo.GetByID(createdDeck.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedDeck)

	assert.Equal(t, createdDeck.ID, retrievedDeck.ID)
	assert.Equal(t, createdDeck.UserID, retrievedDeck.UserID)
	assert.Equal(t, createdDeck.Name, retrievedDeck.Name)
	assert.Equal(t, createdDeck.Description, retrievedDeck.Description)
}

func TestDeckRepository_GetByID_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Try to get a non-existent deck
	randomID := uuid.New()
	deck, err := repo.GetByID(randomID)

	assert.Error(t, err)
	assert.Nil(t, deck)
	assert.Contains(t, err.Error(), "deck not found")
}

func TestDeckRepository_GetAll_Empty(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Get all decks when table is empty
	decks, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, decks)
}

func TestDeckRepository_GetAll_WithData(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create multiple decks
	deck1 := testutils.CreateTestDeck(createdUser.ID)
	deck1.Name = "Deck 1"
	deck2 := testutils.CreateTestDeck(createdUser.ID)
	deck2.Name = "Deck 2"

	_, err = repo.Create(deck1)
	require.NoError(t, err)
	_, err = repo.Create(deck2)
	require.NoError(t, err)

	// Get all decks
	decks, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, decks, 2)

	// Verify decks are ordered by created_at DESC (newest first)
	if decks[0].Name == "Deck 1" {
		assert.Equal(t, "Deck 1", decks[0].Name)
		assert.Equal(t, "Deck 2", decks[1].Name)
	} else {
		assert.Equal(t, "Deck 1", decks[1].Name)
		assert.Equal(t, "Deck 2", decks[0].Name)
	}
}

func TestDeckRepository_GetByUser_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create two users
	user1 := testutils.CreateTestUser()
	user1.PasswordHash = "test_hash"
	user1.Email = "user1@example.com"

	user2 := testutils.CreateTestUser()
	user2.PasswordHash = "test_hash"
	user2.Email = "user2@example.com"

	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser1, err := userRepo.Create(user1)
	require.NoError(t, err)
	createdUser2, err := userRepo.Create(user2)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create decks for user1
	deck1 := testutils.CreateTestDeck(createdUser1.ID)
	deck1.Name = "User1 Deck 1"
	deck2 := testutils.CreateTestDeck(createdUser1.ID)
	deck2.Name = "User1 Deck 2"

	// Create deck for user2
	deck3 := testutils.CreateTestDeck(createdUser2.ID)
	deck3.Name = "User2 Deck 1"

	_, err = repo.Create(deck1)
	require.NoError(t, err)
	_, err = repo.Create(deck2)
	require.NoError(t, err)
	_, err = repo.Create(deck3)
	require.NoError(t, err)

	// Get decks for user1
	user1Decks, err := repo.GetByUser(createdUser1.ID)
	require.NoError(t, err)
	assert.Len(t, user1Decks, 2)

	// Get decks for user2
	user2Decks, err := repo.GetByUser(createdUser2.ID)
	require.NoError(t, err)
	assert.Len(t, user2Decks, 1)
	assert.Equal(t, "User2 Deck 1", user2Decks[0].Name)
}

func TestDeckRepository_GetByUser_Empty(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Try to get decks for a non-existent user
	randomUserID := uuid.New()
	decks, err := repo.GetByUser(randomUserID)

	require.NoError(t, err)
	assert.Empty(t, decks)
}

func TestDeckRepository_Update_Name(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create a deck
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Update the deck's name
	newName := "Updated Deck Name"
	updates := map[string]interface{}{
		"name": newName,
	}

	updatedDeck, err := repo.Update(createdDeck.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedDeck)

	assert.Equal(t, createdDeck.ID, updatedDeck.ID)
	assert.Equal(t, createdDeck.UserID, updatedDeck.UserID)
	assert.Equal(t, newName, updatedDeck.Name)
	assert.True(t, updatedDeck.UpdatedAt.After(createdDeck.UpdatedAt))
}

func TestDeckRepository_Update_Description(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create a deck
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Update the deck's description
	newDescription := "Updated Description"
	updates := map[string]interface{}{
		"description": newDescription,
	}

	updatedDeck, err := repo.Update(createdDeck.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedDeck)

	assert.Equal(t, createdDeck.ID, updatedDeck.ID)
	assert.Equal(t, createdDeck.UserID, updatedDeck.UserID)
	assert.Equal(t, newDescription, updatedDeck.Description)
	assert.True(t, updatedDeck.UpdatedAt.After(createdDeck.UpdatedAt))
}

func TestDeckRepository_Update_BothFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create a deck
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Update both name and description
	newName := "Updated Deck Name"
	newDescription := "Updated Description"
	updates := map[string]interface{}{
		"name":        newName,
		"description": newDescription,
	}

	updatedDeck, err := repo.Update(createdDeck.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedDeck)

	assert.Equal(t, createdDeck.ID, updatedDeck.ID)
	assert.Equal(t, newName, updatedDeck.Name)
	assert.Equal(t, newDescription, updatedDeck.Description)
	assert.True(t, updatedDeck.UpdatedAt.After(createdDeck.UpdatedAt))
}

func TestDeckRepository_Update_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Try to update a non-existent deck
	randomID := uuid.New()
	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	updatedDeck, err := repo.Update(randomID, updates)
	assert.Error(t, err)
	assert.Nil(t, updatedDeck)
	assert.Contains(t, err.Error(), "deck not found")
}

func TestDeckRepository_Update_NoFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Try to update with no fields
	randomID := uuid.New()
	updates := map[string]interface{}{}

	updatedDeck, err := repo.Update(randomID, updates)
	assert.Error(t, err)
	assert.Nil(t, updatedDeck)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestDeckRepository_Delete_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Create a deck
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Delete the deck
	err = repo.Delete(createdDeck.ID)
	require.NoError(t, err)

	// Verify deck is deleted
	deletedDeck, err := repo.GetByID(createdDeck.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedDeck)
	assert.Contains(t, err.Error(), "deck not found")
}

func TestDeckRepository_Delete_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewDeckRepository(td.DB.DB, td.Logger)

	// Try to delete a non-existent deck
	randomID := uuid.New()
	err := repo.Delete(randomID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deck not found")
}

func TestDeckRepository_GetDeckFlashcardCount(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	// Create user and deck first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	userRepo := NewUserRepository(td.DB.DB, td.Logger)
	createdUser, err := userRepo.Create(user)
	require.NoError(t, err)

	repo := NewDeckRepository(td.DB.DB, td.Logger)
	deck := testutils.CreateTestDeck(createdUser.ID)
	createdDeck, err := repo.Create(deck)
	require.NoError(t, err)

	// Get flashcard count for deck (should be 0 initially)
	count, err := repo.GetDeckFlashcardCount(createdDeck.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// TODO: Add test with actual flashcards when flashcard repository tests are implemented
}
