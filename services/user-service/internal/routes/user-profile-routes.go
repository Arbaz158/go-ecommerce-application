package routes

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-ecommerce-application/pkg/middleware"
	"github.com/go-ecommerce-application/services/user-service/internal/handler"
)

func RegisterUserProfileRoutes(r *gin.Engine, h *handler.UserProfileHandler) {
	// Apply JWT middleware to protect all user routes
	userRoutes := r.Group("/users")
	userRoutes.GET("/me", h.GetMe)
	userRoutes.POST("/address", h.CreateAddress)
	userRoutes.GET("/address", h.GetAddresses)
}
