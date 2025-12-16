package dto

// UserResponse represents the user data returned to clients
type UserResponse struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"johndoe"`
}
