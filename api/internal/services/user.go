package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"swipelearn-api/internal/models"
	"swipelearn-api/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
	Logger   *logrus.Logger
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *logrus.Logger) *UserService {
	return &UserService{
		userRepo: repo,
		Logger:   logger,
	}
}

// Create creates a new user with business logic validation
func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		ID:    uuid.New(),
		Email: req.Email,
		Name:  req.Name,
	}

	savedUser, err := s.userRepo.Create(user)
	if err != nil {
		s.Logger.WithError(err).Error("Service failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id": savedUser.ID,
		"email":   savedUser.Email,
	}).Info("User created successfully")

	return savedUser, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", id).Error("Service failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.Logger.WithError(err).WithField("email", email).Error("Service failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAll retrieves all users
func (s *UserService) GetAll() ([]*models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Service failed to get all users")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	s.Logger.WithField("user_count", len(users)).Info("Retrieved all users")
	return users, nil
}

// Update updates a user with business logic validation
func (s *UserService) Update(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user first
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Email != nil && *req.Email != existingUser.Email {
		// Check if new email is already taken
		existingEmailUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingEmailUser != nil && existingEmailUser.ID != id {
			return nil, fmt.Errorf("email %s is already taken", *req.Email)
		}
		updates["email"] = *req.Email
	}

	if req.Name != nil && *req.Name != existingUser.Name {
		updates["name"] = *req.Name
	}

	if len(updates) == 0 {
		return existingUser, nil // No changes needed
	}

	updatedUser, err := s.userRepo.Update(id, updates)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", id).Error("Service failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User updated successfully")

	return updatedUser, nil
}

// Delete removes a user with validation
func (s *UserService) Delete(id uuid.UUID) error {
	// Check if user exists first
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		s.Logger.WithField("user_id", id).Warn("Attempted to delete non-existent user")
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.userRepo.Delete(id)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", id).Error("Service failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.Logger.WithField("user_id", id).Info("User deleted successfully")
	return nil
}
