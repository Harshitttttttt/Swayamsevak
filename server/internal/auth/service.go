package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrEmailInUse         = errors.New("email already in use")
)

// AuthService provides authentication functionality
type AuthService struct {
	userRepo         *models.UserRepository
	refreshTokenRepo *models.RefreshTokenRepository
	jwtSecret        []byte
	accessTokenTTL   time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *models.UserRepository, refreshTokenRepo *models.RefreshTokenRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
		accessTokenTTL:   accessTokenTTL,
	}
}

// Register creates a new user with the provided credentials
func (s *AuthService) Register(email, username, password string) (*models.User, error) {
	// Check if the user already exists
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, ErrEmailInUse
	}

	// Only proceed if the error was "user not found"
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user, err := s.userRepo.CreateUser(email, username, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates the user and returns an access token
func (s *AuthService) Login(email, password string) (string, error) {
	// Get the user from the database
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Verify the password
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Update last_login and updated_at
	// We don't write error handling here which is returned from UpdateLastLogin
	// because we don't want to block user login even if update fails
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generate an access token
	token, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateAccessToken creates a new JWT access token
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	// Set the expiration time
	expirationTime := time.Now().Add(s.accessTokenTTL)

	// Create the JWT Claims
	claims := jwt.MapClaims{
		"sub":      user.ID,               // subject (user ID)
		"username": user.Username,         // custom claim
		"email":    user.Email,            // custom claim
		"exp":      expirationTime.Unix(), // expiration time
		"iat":      time.Now().Unix(),     // issued at time
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// LoginWithRefresh authenticates a user and returns both access and refresh tokens
func (s *AuthService) LoginWithRefresh(email, password string, refreshTokenTTL time.Duration) (accessToken string, refreshToken string, err error) {
	// Get the user from the database
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify the password
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Update last_login and updated_at
	// We don't write error handling here which is returned from UpdateLastLogin
	// because we don't want to block user login even if update fails
	_ = s.userRepo.UpdateLastLogin(user.ID)

	if err := s.refreshTokenRepo.RevokeAllForUser(user.ID); err != nil {
		return "", "", err
	}

	// Generate an access token
	accessToken, err = s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// Create a refresh token
	refeshToken, err := s.refreshTokenRepo.CreateRefreshToken(user.ID, refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refeshToken.Token, nil
}

// RefreshAccessToken creates a new Access Token using a Refresh Token
func (s *AuthService) RefreshAccessToken(oldRefreshTokenString string, refreshTTL time.Duration) (string, string, error) {
	// Rotate the refresh token
	newToken, err := s.refreshTokenRepo.RotateRefreshToken(oldRefreshTokenString, refreshTTL)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Get the user
	user, err := s.userRepo.GetUserByID(newToken.UserID)
	if err != nil {
		return "", "", err
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, newToken.Token, nil
}

func (s *AuthService) Logout(refreshTokenString string) error {
	token, err := s.refreshTokenRepo.GetRefreshToken(refreshTokenString)
	if err != nil {
		return ErrInvalidToken
	}

	if token.Revoked {
		return nil
	}

	return s.refreshTokenRepo.RevokeRefreshToken(refreshTokenString)
}
