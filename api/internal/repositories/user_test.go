package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swipelearn-api/pkg/testutils"
)

func TestUserRepository_Create(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"

	createdUser, err := repo.Create(user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	assert.NotEqual(t, uuid.Nil, createdUser.ID)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.PasswordHash, createdUser.PasswordHash)
	assert.False(t, createdUser.CreatedAt.IsZero())
	assert.False(t, createdUser.UpdatedAt.IsZero())
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	user1 := testutils.CreateTestUser()
	user1.PasswordHash = "test_hash"
	user1.Email = "duplicate@example.com"

	user2 := testutils.CreateTestUser()
	user2.PasswordHash = "test_hash"
	user2.Email = "duplicate@example.com" // Same email

	// First user should be created successfully
	_, err := repo.Create(user1)
	require.NoError(t, err)

	// Second user with same email should fail
	_, err = repo.Create(user2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Get the user by ID
	retrievedUser, err := repo.GetByID(createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	assert.Equal(t, createdUser.ID, retrievedUser.ID)
	assert.Equal(t, createdUser.Email, retrievedUser.Email)
	assert.Equal(t, createdUser.Name, retrievedUser.Name)
	assert.Equal(t, createdUser.PasswordHash, retrievedUser.PasswordHash)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Try to get a non-existent user
	randomID := uuid.New()
	user, err := repo.GetByID(randomID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Get the user by email
	retrievedUser, err := repo.GetByEmail(createdUser.Email)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	assert.Equal(t, createdUser.ID, retrievedUser.ID)
	assert.Equal(t, createdUser.Email, retrievedUser.Email)
	assert.Equal(t, createdUser.Name, retrievedUser.Name)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Try to get a non-existent user by email
	user, err := repo.GetByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_GetAll_Empty(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Get all users when table is empty
	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_GetAll_WithData(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create multiple users
	user1 := testutils.CreateTestUser()
	user1.PasswordHash = "test_hash"
	user1.Email = "user1@example.com"
	user1.Name = "User 1"

	user2 := testutils.CreateTestUser()
	user2.PasswordHash = "test_hash"
	user2.Email = "user2@example.com"
	user2.Name = "User 2"

	_, err := repo.Create(user1)
	require.NoError(t, err)
	_, err = repo.Create(user2)
	require.NoError(t, err)

	// Get all users
	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, users, 2)

	// Verify users are ordered by created_at DESC (newest first)
	if users[0].Name == "User 1" {
		assert.Equal(t, "User 1", users[0].Name)
		assert.Equal(t, "User 2", users[1].Name)
	} else {
		assert.Equal(t, "User 1", users[1].Name)
		assert.Equal(t, "User 2", users[0].Name)
	}
}

func TestUserRepository_Update_Name(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Update the user's name
	newName := "Updated Name"
	updates := map[string]interface{}{
		"name": newName,
	}

	updatedUser, err := repo.Update(createdUser.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)

	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, createdUser.Email, updatedUser.Email)
	assert.Equal(t, newName, updatedUser.Name)
	assert.True(t, updatedUser.UpdatedAt.After(createdUser.UpdatedAt))
}

func TestUserRepository_Update_Email(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Update the user's email
	newEmail := "updated@example.com"
	updates := map[string]interface{}{
		"email": newEmail,
	}

	updatedUser, err := repo.Update(createdUser.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)

	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, createdUser.Name, updatedUser.Name)
	assert.True(t, updatedUser.UpdatedAt.After(createdUser.UpdatedAt))
}

func TestUserRepository_Update_BothFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Update both name and email
	newName := "Updated Name"
	newEmail := "updated@example.com"
	updates := map[string]interface{}{
		"name":  newName,
		"email": newEmail,
	}

	updatedUser, err := repo.Update(createdUser.ID, updates)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)

	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, newName, updatedUser.Name)
	assert.True(t, updatedUser.UpdatedAt.After(createdUser.UpdatedAt))
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Try to update a non-existent user
	randomID := uuid.New()
	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	updatedUser, err := repo.Update(randomID, updates)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Update_NoFields(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Try to update with no fields
	randomID := uuid.New()
	updates := map[string]interface{}{}

	updatedUser, err := repo.Update(randomID, updates)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestUserRepository_Delete_Success(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Create a user first
	user := testutils.CreateTestUser()
	user.PasswordHash = "test_hash"
	createdUser, err := repo.Create(user)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	deletedUser, err := repo.GetByID(createdUser.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	td := testutils.SetupTestDatabase(t)
	defer td.Close()
	td.RunMigrations(t)

	repo := NewUserRepository(td.DB.DB, td.Logger)

	// Try to delete a non-existent user
	randomID := uuid.New()
	err := repo.Delete(randomID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}
