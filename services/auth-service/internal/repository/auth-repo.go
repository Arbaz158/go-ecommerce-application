package repository

import (
	"errors"
	"log"

	"github.com/go-ecommerce-application/services/auth-service/internal/database"
	"github.com/go-ecommerce-application/services/auth-service/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user models.AuthUser) error
	GetUserByEmail(email string) (*models.AuthUser, error)
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (r *authRepository) CreateUser(user models.AuthUser) error {

	err := database.DB.Create(&user).Error
	if err != nil {
		log.Println("Error while creating user :", err)
		return err
	}
	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.AuthUser, error) {
	var authUser models.AuthUser

	err := database.DB.
		Where("email = ?", email).
		First(&authUser).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Println("Error while getting user data:", err)
		return nil, err
	}

	return &authUser, nil
}
