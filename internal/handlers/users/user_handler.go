package users

import (
	"net/http"

	userService "farm-backend/internal/services/users"
	"farm-backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *userService.UserService
}

func NewUserHandler(svc *userService.UserService) *UserHandler {
	return &UserHandler{UserService: svc}
}

// GetProfile retrieves the currently logged-in user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(uint)

	user, err := h.UserService.GetProfile(userID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": user})
}

// UpdateProfile updates the currently logged-in user's profile details
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(uint)

	var req validation.UpdateProfileRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}

	if err := h.UserService.UpdateProfile(userID, &req); err != nil {
		validation.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ChangePassword changes the currently logged-in user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(uint)

	var req validation.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}

	if err := h.UserService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		validation.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}