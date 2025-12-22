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

// RotateRefreshToken atomically creates a new refresh token for the same user and revokes the old one.
// Returns the new RefreshToken or an error (sql.ErrNoRows if old token not found).
func (r *RefreshTokenRepository) RotateRefreshToken(oldTokenString string, ttl time.Duration) (*RefreshToken, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	// on failure rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var userID uuid.UUID
	var revoked bool
	var expiresAt time.Time

	// Retrieve the old token's user ID
	query := `
		SELECT user_id, revoked, expires_at
		FROM refresh_tokens
		WHERE token = $1
		FOR UPDATE;
	`

	err = tx.QueryRow(query, oldTokenString).Scan(&userID, &revoked, &expiresAt)
	if err != nil {
		return nil, err
	}

	// Create a new refresh token
	newTokenID := uuid.New()
	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	newTokenString := base64.RawURLEncoding.EncodeToString(tokenBytes)
	newExpiresAt := time.Now().Add(ttl)

	insertQuery :=
		`
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at, revoked;
	`

	var newToken RefreshToken
	newToken.ID = newTokenID
	newToken.UserID = userID
	newToken.Token = newTokenString
	newToken.ExpiresAt = newExpiresAt

	if err = tx.QueryRow(insertQuery, newToken.ID, newToken.UserID, newToken.Token, newToken.ExpiresAt).Scan(&newToken.CreatedAt, &newToken.UpdatedAt, &newToken.Revoked); err != nil {
		return nil, err
	}

	// Revoke the old token
	updateQuery :=
		`
		UPDATE refresh_tokens
		SET revoked = true, updated_at = NOW()
		WHERE token = $1;
	`

	if _, err = tx.Exec(updateQuery, oldTokenString); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &newToken, nil
}

func (r *RefreshTokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	query :=
		`
		UPDATE refreh_tokens
		SET revoked = true, updated_at = NOW()
		WHERE user_id = $1 AND revoked = false;
	`

	_, err := r.db.Exec(query, userID)

	return err
}
