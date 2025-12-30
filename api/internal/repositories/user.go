package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"swipelearn-api/internal/models"
)

type UserRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewUserRepository(db *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		DB:     db,
		Logger: logger,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, email, name, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, password_hash, created_at, updated_at
	`

	err := r.DB.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		r.Logger.WithError(err).Error("Failed to create user in database")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.Logger.WithField("user_id", user.ID).Info("User created successfully")
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.Logger.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.Logger.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to get all users")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			r.Logger.WithError(err).Error("Failed to scan user row")
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		r.Logger.WithError(err).Error("Error after scanning user rows")
		return nil, fmt.Errorf("error scanning users: %w", err)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	// Build dynamic UPDATE query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if name, ok := updates["name"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if email, ok := updates["email"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, email)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add updated_at and id
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, email, name, password_hash, created_at, updated_at
	`, fmt.Sprintf("%s", setParts), argIndex)

	user := &models.User{}
	err := r.DB.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.Logger.WithError(err).WithField("user_id", id).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	r.Logger.WithField("user_id", user.ID).Info("User updated successfully")
	return user, nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		r.Logger.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.Logger.WithField("user_id", id).Info("User deleted successfully")
	return nil
}
