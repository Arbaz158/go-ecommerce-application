package repository

import (
	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetUserProfileByID(id uint) (*models.UserProfile, error)
	SaveAddress(address *models.Address) error
	GetUserAddresses(userID uint) ([]models.Address, error)
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (r *userProfileRepository) GetUserProfileByID(id uint) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	if err := r.db.First(&userProfile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userProfile, nil
}

func (r *userProfileRepository) SaveAddress(address *models.Address) error {
	return r.db.Create(address).Error
}

func (r *userProfileRepository) GetUserAddresses(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	if err := r.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}
