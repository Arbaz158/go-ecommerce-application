package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/services/auth-service/internal/dto"
	"github.com/go-ecommerce-application/services/auth-service/internal/models"
	"github.com/go-ecommerce-application/services/auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignupHandler(c *gin.Context) {
	var authData models.AuthUser
	if err := c.ShouldBindJSON(&authData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err := h.authService.Signup(authData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Signup failed"})
		return
	}
	c.JSON(201, gin.H{"message": "Signup successful"})
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var loginData dto.LoginRequest
	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Request",
		})
		return
	}
	loginResponse, err := h.authService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error occured while login " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"response": loginResponse,
	})
}

func (h *AuthHandler) RefreshHandler(c *gin.Context) {
	// Token refresh logic here
}

func (h *AuthHandler) ProtectedHandler(c *gin.Context) {
	// Protected route logic here
}
