package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RefreshTokenRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewRefreshTokenRepository(db *sql.DB, logger *logrus.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB:     db,
		Logger: logger,
	}
}

// StoreRefreshToken stores a hashed refresh token in database
func (r *RefreshTokenRepository) StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	// Hash token for storage using SHA256 (more suitable for long strings)
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.DB.Exec(
		query,
		uuid.New(), // Generate new ID for token record
		userID,
		tokenHash,
		expiresAt,
	)

	if err != nil {
		r.Logger.WithError(err).Error("Failed to store refresh token")
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

// GetValidRefreshToken retrieves a valid (non-expired, non-revoked) refresh token
func (r *RefreshTokenRepository) GetValidRefreshToken(userID uuid.UUID, tokenString string) (*RefreshToken, error) {
	// Get all unexpired tokens for user
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to query refresh tokens")
		return nil, fmt.Errorf("failed to query refresh tokens: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		token := &RefreshToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.TokenHash,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.RevokedAt,
		)
		if err != nil {
			r.Logger.WithError(err).Error("Failed to scan refresh token")
			continue
		}

		// Check if the provided token matches the stored hash
		computedHash := sha256.Sum256([]byte(tokenString))
		storedHash, _ := hex.DecodeString(token.TokenHash)
		if string(computedHash[:]) == string(storedHash[:]) {
			return token, nil
		}
	}

	return nil, fmt.Errorf("valid refresh token not found")
}

// RevokeToken revokes a refresh token
func (r *RefreshTokenRepository) RevokeToken(tokenID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`

	_, err := r.DB.Exec(query, tokenID)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to revoke refresh token")
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// RevokeUserTokens revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeUserTokens(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`

	_, err := r.DB.Exec(query, userID)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to revoke user refresh tokens")
		return fmt.Errorf("failed to revoke user refresh tokens: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired refresh tokens
func (r *RefreshTokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.DB.Exec(query)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to cleanup expired tokens")
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
