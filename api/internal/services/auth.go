package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"swipelearn-api/internal/models"
	"swipelearn-api/internal/repositories"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	jwtService       *JWTService
	Logger           *logrus.Logger
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	jwtService *JWTService,
	logger *logrus.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		Logger:           logger,
	}
}

// Register creates a new user with password
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := s.jwtService.HashPassword(req.Password)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to process password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	savedUser, err := s.userRepo.Create(user)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id": savedUser.ID,
		"email":   savedUser.Email,
	}).Info("User registered successfully")

	return savedUser, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		s.Logger.WithField("email", req.Email).Warn("Login attempt with non-existent email")
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if !s.jwtService.CheckPassword(req.Password, user.PasswordHash) {
		s.Logger.WithField("email", req.Email).Warn("Login attempt with invalid password")
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to generate tokens")
	}

	// Store refresh token
	err = s.refreshTokenRepo.StoreRefreshToken(
		user.ID,
		refreshToken,
		time.Now().Add(s.jwtService.refreshTokenTTL),
	)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to store refresh token")
		return nil, fmt.Errorf("failed to store refresh token")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	s.Logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates new tokens from a valid refresh token
func (s *AuthService) RefreshToken(req *models.RefreshRequest) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		s.Logger.WithError(err).Warn("Invalid refresh token")
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user to ensure they still exist
	user, err := s.userRepo.GetByID(uuid.MustParse(claims.UserID))
	if err != nil {
		s.Logger.WithError(err).Error("User not found during refresh")
		return nil, fmt.Errorf("user not found")
	}

	// Verify refresh token exists in database
	_, err = s.refreshTokenRepo.GetValidRefreshToken(uuid.MustParse(claims.UserID), req.RefreshToken)
	if err != nil {
		s.Logger.WithError(err).Warn("Refresh token not found in database")
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Revoke the old refresh token
	err = s.refreshTokenRepo.RevokeUserTokens(uuid.MustParse(claims.UserID))
	if err != nil {
		s.Logger.WithError(err).Warn("Failed to revoke old refresh tokens")
		// Continue anyway - this is not fatal
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to generate new tokens")
		return nil, fmt.Errorf("failed to generate tokens")
	}

	// Store new refresh token
	err = s.refreshTokenRepo.StoreRefreshToken(
		user.ID,
		newRefreshToken,
		time.Now().Add(s.jwtService.refreshTokenTTL),
	)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to store new refresh token")
		return nil, fmt.Errorf("failed to store refresh token")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	s.Logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("Token refreshed successfully")

	return &models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// Logout revokes all refresh tokens for a user
func (s *AuthService) Logout(userID uuid.UUID) error {
	err := s.refreshTokenRepo.RevokeUserTokens(userID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to revoke tokens during logout")
		return fmt.Errorf("failed to logout")
	}

	s.Logger.WithField("user_id", userID).Info("User logged out successfully")
	return nil
}
