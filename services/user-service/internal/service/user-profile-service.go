package service

import (
	"log"

	"github.com/go-ecommerce-application/pkg/kafka/events"
	"github.com/go-ecommerce-application/services/user-service/internal/models"
	"github.com/go-ecommerce-application/services/user-service/internal/repository"
)

type UserProfileService interface {
	// GetUserProfile(id int) (*models.UserProfile, error)
	GetUserProfileByUserID(userID string) (*models.UserProfile, error)
	HandleUserSignedUpEvent(event *events.UserSignedUp) error
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

func (u *userProfileService) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	userProfile, err := u.userProfileRepo.GetUserProfileByUserID(userID)
	if err != nil {
		log.Println("error while fetching user profile by user id :", err)
		return nil, err
	}
	return userProfile, nil
}

func (u *userProfileService) HandleUserSignedUpEvent(event *events.UserSignedUp) error {
	// Check if user profile already exists
	existing, err := u.userProfileRepo.GetUserProfileByUserID(event.UserID)
	if err != nil {
		log.Println("error while checking existing user profile :", err)
		return err
	}
	if existing != nil {
		log.Printf("user profile already exists for user id %s", event.UserID)
		return nil
	}

	// Create new user profile
	profile := &models.UserProfile{
		UserID:    event.UserID,
		Email:     event.Email,
		FirstName: event.FirstName,
		LastName:  event.LastName,
	}

	err = u.userProfileRepo.CreateUserProfile(profile)
	if err != nil {
		log.Println("error while creating user profile :", err)
		return err
	}

	log.Printf("user profile created successfully for user id %s", event.UserID)
	return nil
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
