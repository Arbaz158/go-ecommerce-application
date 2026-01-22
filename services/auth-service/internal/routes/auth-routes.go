package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/services/auth-service/internal/handler"
)

func RegisterAuthRoutes(r *gin.Engine, h *handler.AuthHandler) {
	auth := r.Group("/auth")

	auth.POST("/signup", h.SignupHandler)
	auth.POST("/login", h.LoginHandler)
	auth.POST("/refresh", h.RefreshHandler)

	auth.GET("/logout", h.ProtectedHandler)
}
