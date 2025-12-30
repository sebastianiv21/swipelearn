package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swipelearn-api/pkg/testutils"
)

func TestJWTService_HashPassword(t *testing.T) {
	logger := testutils.TestLogger()

	// Set a secret for testing
	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	password := "test_password_123"
	hash, err := service.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash) // Hash should be different from password
}

func TestJWTService_CheckPassword_Correct(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	password := "test_password_123"
	hash, err := service.HashPassword(password)
	require.NoError(t, err)

	// Check password against its hash
	isValid := service.CheckPassword(password, hash)
	assert.True(t, isValid)
}

func TestJWTService_CheckPassword_Incorrect(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	password := "test_password_123"
	wrongPassword := "wrong_password"
	hash, err := service.HashPassword(password)
	require.NoError(t, err)

	// Check wrong password against hash
	isValid := service.CheckPassword(wrongPassword, hash)
	assert.False(t, isValid)
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	userID := uuid.New().String()
	email := "test@example.com"

	accessToken, refreshToken, err := service.GenerateTokenPair(userID, email)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, accessToken, refreshToken) // Should be different tokens
}

func TestJWTService_ValidateAccessToken_Valid(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	userID := uuid.New().String()
	email := "test@example.com"

	// Generate token first
	accessToken, _, err := service.GenerateTokenPair(userID, email)
	require.NoError(t, err)

	// Validate the token
	claims, err := service.ValidateAccessToken(accessToken)

	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims content
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWTService_ValidateAccessToken_Invalid(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	invalidToken := "invalid.jwt.token"

	// Validate invalid token
	claims, err := service.ValidateAccessToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateRefreshToken_Valid(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	userID := uuid.New().String()

	// Generate refresh token first
	_, refreshToken, err := service.GenerateTokenPair(userID, "test@example.com")
	require.NoError(t, err)

	// Validate refresh token
	claims, err := service.ValidateRefreshToken(refreshToken)

	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims content
	assert.Equal(t, userID, claims.UserID)
	assert.NotNil(t, claims.TokenID) // Refresh tokens should have token ID
}

func TestJWTService_ValidateRefreshToken_Invalid(t *testing.T) {
	logger := testutils.TestLogger()

	t.Setenv("JWT_SECRET", "test_secret_key")
	service := NewJWTService(logger)

	invalidToken := "invalid.jwt.token"

	// Validate invalid refresh token
	claims, err := service.ValidateRefreshToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
