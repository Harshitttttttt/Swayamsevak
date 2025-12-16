package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

// Refresh Token represents	a refresh token in the system
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Revoked   bool
}

// RefreshTokenRepository handles database operations for refresh tokens
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new refresh token repo
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// CreateRefreshToken creates a new refresh token for a user
func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, ttl time.Duration) (*RefreshToken, error) {
	// Generate a unique token identifier
	tokenID := uuid.New()
	// tokenString := tokenID.String()
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	tokenString := base64.RawURLEncoding.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(ttl)

	token := &RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	query :=
		`
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at, revoked;
	`

	err := r.db.QueryRow(query, token.ID, token.UserID, token.Token, token.ExpiresAt).Scan(&token.CreatedAt, &token.UpdatedAt, &token.Revoked)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetRefreshToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	query :=
		`
		SELECT id, user_id, token, expires_at, created_at, updated_at, revoked
		FROM refresh_tokens
		WHERE token = $1;
	`

	var token RefreshToken
	err := r.db.QueryRow(query, tokenString).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt, &token.Revoked)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// RevokeRefreshToken marks a RefreshToken as revoked
func (r *RefreshTokenRepository) RevokeRefreshToken(tokenString string) error {
	query :=
		`
		UPDATE refresh_tokens
		SET revoked = true, updated_at = NOW()
		WHERE token = $1;
	`

	_, err := r.db.Exec(query, tokenString)
	return err
}
