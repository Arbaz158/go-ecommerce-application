package repository

import (
	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetUserProfileByUserID(userID string) (*models.UserProfile, error)
	CreateUserProfile(profile *models.UserProfile) error
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

func (r *userProfileRepository) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	if err := r.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userProfile, nil
}

func (r *userProfileRepository) CreateUserProfile(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
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
