package service

import (
	"errors"
	"testing"

	"github.com/go-ecommerce-application/services/user-service/internal/domain/models"
)

// MockUserProfileRepository is a mock implementation of UserProfileRepository
type MockUserProfileRepository struct {
	GetUserProfileByUserIDFunc func(userID string) (*models.UserProfile, error)
	CreateUserProfileFunc      func(profile *models.UserProfile) error
	SaveAddressFunc            func(address *models.Address) error
	GetUserAddressesFunc       func(userID uint) ([]models.Address, error)
}

// Implement repository methods for the mock

func (m *MockUserProfileRepository) CreateUserProfile(profile *models.UserProfile) error {
	if m.CreateUserProfileFunc != nil {
		return m.CreateUserProfileFunc(profile)
	}
	return nil
}

func (m *MockUserProfileRepository) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	if m.GetUserProfileByUserIDFunc != nil {
		return m.GetUserProfileByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *MockUserProfileRepository) SaveAddress(address *models.Address) error {
	if m.SaveAddressFunc != nil {
		return m.SaveAddressFunc(address)
	}
	return nil
}

func (m *MockUserProfileRepository) GetUserAddresses(userID uint) ([]models.Address, error) {
	if m.GetUserAddressesFunc != nil {
		return m.GetUserAddressesFunc(userID)
	}
	return nil, nil
}

func TestGetUserProfile_Success(t *testing.T) {
	// Arrange - Setup
	expectedProfile := &models.UserProfile{
		ID:        1,
		UserID:    "ewuhiwj23200",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockRepo := &MockUserProfileRepository{
		GetUserProfileByUserIDFunc: func(userID string) (*models.UserProfile, error) {
			if userID == "ewuhiwj23200" {
				return expectedProfile, nil
			}
			return nil, nil
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	profile, err := service.GetUserProfileByUserID("ewuhiwj23200")

	// Assert - Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile == nil {
		t.Fatalf("expected profile, got nil")
	}

	if profile.FirstName != "John" || profile.LastName != "Doe" {
		t.Errorf("expected name 'John Doe', got %s %s", profile.FirstName, profile.LastName)
	}

	if profile.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got %s", profile.Email)
	}
}

func TestGetUserProfile_NotFound(t *testing.T) {
	// Arrange - Setup
	mockRepo := &MockUserProfileRepository{
		GetUserProfileByUserIDFunc: func(userID string) (*models.UserProfile, error) {
			return nil, nil // User not found returns nil profile
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	profile, err := service.GetUserProfileByUserID("ewuhiwj23200")

	// Assert - Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile != nil {
		t.Fatalf("expected nil profile, got %v", profile)
	}
}

func TestGetUserProfile_RepositoryError(t *testing.T) {
	// Arrange - Setup
	expectedError := errors.New("database connection failed")
	mockRepo := &MockUserProfileRepository{
		GetUserProfileByUserIDFunc: func(userID string) (*models.UserProfile, error) {
			return nil, expectedError
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	profile, err := service.GetUserProfileByUserID("ewuhiwj23200")

	// Assert - Verify
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if profile != nil {
		t.Fatalf("expected nil profile on error, got %v", profile)
	}

	if err.Error() != "database connection failed" {
		t.Errorf("expected error message 'database connection failed', got %s", err.Error())
	}
}

func TestSaveAddress_Success(t *testing.T) {
	// Arrange - Setup
	address := models.Address{
		UserID:     1,
		Street:     "123 Main St",
		City:       "Cityville",
		State:      "Stateville",
		PostalCode: "12345",
	}

	mockRepo := &MockUserProfileRepository{
		SaveAddressFunc: func(addr *models.Address) error {
			if addr.UserID == 1 && addr.Street == "123 Main St" {
				return nil
			}
			return errors.New("invalid address")
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	err := service.SaveAddress(address)

	// Assert - Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSaveAddress_Error(t *testing.T) {
	// Arrange - Setup
	address := models.Address{
		UserID:     1,
		Street:     "123 Main St",
		City:       "Cityville",
		State:      "Stateville",
		PostalCode: "12345",
	}

	expectedError := errors.New("failed to save address")
	mockRepo := &MockUserProfileRepository{
		SaveAddressFunc: func(addr *models.Address) error {
			return expectedError
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	err := service.SaveAddress(address)

	// Assert - Verify
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "failed to save address" {
		t.Errorf("expected error message 'failed to save address', got %s", err.Error())
	}
}

func TestGetUserAddresses_Success(t *testing.T) {
	// Arrange - Setup
	expectedAddresses := []models.Address{
		{
			ID:         1,
			UserID:     1,
			Street:     "123 Main St",
			City:       "Cityville",
			State:      "Stateville",
			PostalCode: "12345",
		},
		{
			ID:         2,
			UserID:     1,
			Street:     "456 Side St",
			City:       "Townsville",
			State:      "Regionville",
			PostalCode: "67890",
		},
	}

	mockRepo := &MockUserProfileRepository{
		GetUserAddressesFunc: func(userID uint) ([]models.Address, error) {
			if userID == 1 {
				return expectedAddresses, nil
			}
			return []models.Address{}, nil
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	addresses, err := service.GetUserAdresses(1)

	// Assert - Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(addresses) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(addresses))
	}

	if addresses[0].Street != "123 Main St" {
		t.Errorf("expected street '123 Main St', got %s", addresses[0].Street)
	}
}

func TestGetUserAddresses_Empty(t *testing.T) {
	// Arrange - Setup
	mockRepo := &MockUserProfileRepository{
		GetUserAddressesFunc: func(userID uint) ([]models.Address, error) {
			return []models.Address{}, nil
		},
	}

	service := NewUserProfileService(mockRepo)

	// Act - Execute
	addresses, err := service.GetUserAdresses(1)

	// Assert - Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(addresses) != 0 {
		t.Errorf("expected 0 addresses, got %d", len(addresses))
	}
}
