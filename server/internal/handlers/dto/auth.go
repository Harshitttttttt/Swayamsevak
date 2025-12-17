package dto

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"SecurePass123!"`
}

// RegisterResponse contains the user data after successful registration
type RegisterResponse struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"johndoe"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"SecurePass123!"`
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshRequest represents the refresh token payload
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4"`
}

// RefreshResponse contains the new access token
type RefreshResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
