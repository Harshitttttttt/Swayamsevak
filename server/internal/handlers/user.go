package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/handlers/dto"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/middleware"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
)

// UserHandler contains HTTP handlers for user-related endpoints
type UserHandler struct {
	userRepo *models.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *models.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// Profile godoc
// @Summary      Get user profile
// @Description  Retrieve the authenticated user's profile information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.UserResponse
// @Failure      401 {string} string "Unauthorized"
// @Router       /profile [get]
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from the database
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusUnauthorized)
		return
	}

	// Return User Profile
	response := dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
