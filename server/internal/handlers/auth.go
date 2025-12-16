package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/auth"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/handlers/dto"
)

// AuthHandler contains HTTP handlers for authentication
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new Auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email, username, and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration credentials"
// @Success      201 {object} dto.RegisterResponse
// @Failure      400 {string} string "Invalid Request Payload"
// @Failure      409 {string} string "Email already in use"
// @Failure      500 {string} string "Error Creating User"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	// Validate Input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Email, Username and Password are required", http.StatusBadRequest)
		return
	}

	// Call the AuthService to register the user
	user, err := h.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Error Creating User", http.StatusInternalServerError)
		return
	}

	// Return the created user (without sensistive data)
	response := dto.RegisterResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and receive JWT access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body dto.LoginRequest true "Login credentials"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {string} string "Invalid Request Payload"
// @Failure      401 {string} string "Invalid Credentials"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	// Attempt to login
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return the token
	response := dto.LoginResponse{
		Token: token,
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        token body dto.RefreshRequest true "Refresh token"
// @Success      200 {object} dto.RefreshResponse
// @Failure      400 {string} string "Invalid Request Payload"
// @Failure      401 {string} string "Invalid or Expired refresh token"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	// Attempt to refresh the token
	token, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			http.Error(w, "Invalid or Expired refresh token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return the new access token
	response := &dto.RefreshResponse{
		Token: token,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
