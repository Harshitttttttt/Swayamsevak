package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents a user in our system
type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLogin    sql.NullTime
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(email, username, passwordHash string) (*User, error) {
	user := &User{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	}

	query :=
		`
    INSERT INTO users (email, username, password_hash)
    VALUES ($1, $2, $3)
    RETURNING id, created_at, updated_at;
  `

	err := r.db.QueryRow(query, user.Email, user.Username, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retreives a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query :=
		`
		SELECT id, email, username, password_hash, created_at, updated_at, last_login
		FROM users
		WHERE email = $1;
	`

	var user User

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retreives user by their ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*User, error) {
	query :=
		`
		SELECT id, email, username, password_hash, created_at, updated_at, last_login
		FROM users
		WHERE id = $1; 
	`

	var user User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastLogin updates a user's last_login
func (r *UserRepository) UpdateLastLogin(userId uuid.UUID) error {
	query :=
		`
		UPDATE users
		SET last_login = now(), updated_at = now()
		WHERE id = $1;
	`

	_, err := r.db.Exec(query, userId)
	return err
}
