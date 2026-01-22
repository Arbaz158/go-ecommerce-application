package service

import (
	"log"

	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"github.com/go-ecommerce-application/services/user-service/internal/repository"
)

type UserProfileService interface {
	GetUserProfile(id int) (*models.UserProfile, error)
	SaveAddress(adress models.Address) error
	GetUserAdresses(userId int) ([]models.Address, error)
}

type userProfileService struct {
	userProfileRepo repository.UserProfileRepository
}

func NewUserProfileService(userProfileRepository repository.UserProfileRepository) UserProfileService {
	return &userProfileService{
		userProfileRepo: userProfileRepository,
	}
}

func (u *userProfileService) GetUserProfile(id int) (*models.UserProfile, error) {
	// var userProfile models.UserProfile
	userProfile, err := u.userProfileRepo.GetUserProfileByID(uint(id))
	if err != nil {
		log.Println("error while fetching user profile :", err)
		return nil, err
	}
	return userProfile, nil
}

func (u *userProfileService) SaveAddress(adress models.Address) error {
	err := u.userProfileRepo.SaveAddress(&adress)
	if err != nil {
		log.Println("error while saving the adress :", err)
		return err
	}
	return nil
}

func (u *userProfileService) GetUserAdresses(userId int) ([]models.Address, error) {
	addresses, err := u.userProfileRepo.GetUserAddresses(uint(userId))
	if err != nil {
		log.Println("error while fetching user addresses :", err)
		return nil, err
	}
	return addresses, nil
}
