package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/libs/auth"
	"github.com/go-ecommerce-application/services/user-service/internal/handler"
)

func RegisterUserProfileRoutes(r *gin.Engine, h *handler.UserProfileHandler) {
	// Health check - no auth required
	r.GET("/health", h.HealthCheck)

	// Apply JWT middleware to protect all user routes
	userRoutes := r.Group("/users")
	userRoutes.Use(auth.AuthMiddleware())
	{
		userRoutes.GET("/me", h.GetMe)
		userRoutes.POST("/address", h.CreateAddress)
		userRoutes.GET("/address", h.GetAddresses)
	}
}
