package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"github.com/go-ecommerce-application/services/user-service/internal/service"
)

type UserProfileHandler struct {
	userProfileService service.UserProfileService
}

func NewUserProfileHandler(userProfileService service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService: userProfileService,
	}
}

func (h *UserProfileHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "User Profile Service is healthy"})
}

func (h *UserProfileHandler) GetMe(c *gin.Context) {
	userProfile, err := h.userProfileService.GetUserProfile(1)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user profile"})
		return
	}
	c.JSON(200, userProfile)
}

func (h *UserProfileHandler) CreateAddress(c *gin.Context) {
	var userAdresses models.Address
	if err := c.ShouldBindJSON(&userAdresses); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.userProfileService.SaveAddress(userAdresses)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save address"})
		return
	}

	c.JSON(201, gin.H{"message": "Address created successfully"})

}

func (h *UserProfileHandler) GetAddresses(c *gin.Context) {
	userAdress, err := h.userProfileService.GetUserAdresses(2)
	if err != nil {
		c.JSON(500, gin.H{"error": "error while getting user adresses :" + err.Error()})
		return
	}
	c.JSON(200, userAdress)

}
