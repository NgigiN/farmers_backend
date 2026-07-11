package auth

import (
	authService "farm-backend/internal/services/auth"
	"farm-backend/internal/validation"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthService *authService.Service
}

func NewAuthHandler(svc *authService.Service) *Handler {
	return &Handler{AuthService: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req validation.RegisterRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	user := validation.UserFromRegisterRequest(&req)
	if err := h.AuthService.Register(user); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	resp, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	resp, err := h.AuthService.GoogleLogin(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}