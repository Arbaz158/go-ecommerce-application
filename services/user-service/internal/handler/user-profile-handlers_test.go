package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-ecommerce-application/pkg/kafka/events"
	"github.com/go-ecommerce-application/services/user-service/internal/models"
)

// MockUserProfileService is a mock implementation of UserProfileService
type MockUserProfileService struct {
	GetUserProfileByUserIDFunc  func(userID string) (*models.UserProfile, error)
	GetUserProfileFunc          func(id int) (*models.UserProfile, error)
	HandleUserSignedUpEventFunc func(event *events.UserSignedUp) error
	SaveAddressFunc             func(address models.Address) error
	GetUserAdressesFunc         func(userId int) ([]models.Address, error)
}

// Implement service methods for the mock

func (m *MockUserProfileService) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	if m.GetUserProfileByUserIDFunc != nil {
		return m.GetUserProfileByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *MockUserProfileService) GetUserProfile(id int) (*models.UserProfile, error) {
	if m.GetUserProfileFunc != nil {
		return m.GetUserProfileFunc(id)
	}
	return nil, nil
}

func (m *MockUserProfileService) HandleUserSignedUpEvent(event *events.UserSignedUp) error {
	if m.HandleUserSignedUpEventFunc != nil {
		return m.HandleUserSignedUpEventFunc(event)
	}
	return nil
}

func (m *MockUserProfileService) SaveAddress(address models.Address) error {
	if m.SaveAddressFunc != nil {
		return m.SaveAddressFunc(address)
	}
	return nil
}

func (m *MockUserProfileService) GetUserAdresses(userId int) ([]models.Address, error) {
	if m.GetUserAdressesFunc != nil {
		return m.GetUserAdressesFunc(userId)
	}
	return nil, nil
}

// Test Case 1: HealthCheck - Success
func TestHealthCheck_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{}
	handler := NewUserProfileHandler(mockService)

	// Act
	handler.HealthCheck(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "User Profile Service is healthy" {
		t.Errorf("expected status message, got %v", response["status"])
	}
}

// Test Case 2: GetMe - Success
func TestGetMe_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedProfile := &models.UserProfile{
		ID:        1,
		UserID:    "user123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockService := &MockUserProfileService{
		GetUserProfileByUserIDFunc: func(userID string) (*models.UserProfile, error) {
			if userID == "user123" {
				return expectedProfile, nil
			}
			return nil, nil
		},
	}

	handler := NewUserProfileHandler(mockService)

	// Set userID in context as string
	c.Set("userID", "user123")
	c.Request = httptest.NewRequest("GET", "/users/me", nil)

	// Act
	handler.GetMe(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var response models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.FirstName != "John" || response.LastName != "Doe" {
		t.Errorf("expected name 'John Doe', got %s %s", response.FirstName, response.LastName)
	}
}

// Test Case 3: GetMe - Service Error
func TestGetMe_ServiceError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{
		GetUserProfileByUserIDFunc: func(userID string) (*models.UserProfile, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewUserProfileHandler(mockService)
	c.Set("userID", "user123")
	c.Request = httptest.NewRequest("GET", "/users/me", nil)

	// Act
	handler.GetMe(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Failed to fetch user profile" {
		t.Errorf("expected error message, got %v", response["error"])
	}
}

// Test Case 4: GetMe - User Not Found
func TestGetMe_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{
		GetUserProfileFunc: func(id int) (*models.UserProfile, error) {
			return nil, nil
		},
	}

	handler := NewUserProfileHandler(mockService)
	c.Set("userID", 999)
	c.Request = httptest.NewRequest("GET", "/users/me", nil)

	// Act
	handler.GetMe(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var response interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response != nil {
		t.Errorf("expected null profile, got %v", response)
	}
}

// Test Case 5: CreateAddress - Success
func TestCreateAddress_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	address := models.Address{
		UserID:     1,
		Street:     "123 Main St",
		City:       "Cityville",
		State:      "Stateville",
		PostalCode: "12345",
	}

	mockService := &MockUserProfileService{
		SaveAddressFunc: func(addr models.Address) error {
			return nil
		},
	}

	handler := NewUserProfileHandler(mockService)

	// Create JSON body
	body, _ := json.Marshal(address)
	c.Request = httptest.NewRequest("POST", "/users/address", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.CreateAddress(c)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status code 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Address created successfully" {
		t.Errorf("expected success message, got %v", response["message"])
	}
}

// Test Case 6: CreateAddress - Invalid JSON
func TestCreateAddress_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{}
	handler := NewUserProfileHandler(mockService)

	// Create invalid JSON body
	c.Request = httptest.NewRequest("POST", "/users/address", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.CreateAddress(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Invalid request payload" {
		t.Errorf("expected error message 'Invalid request payload', got %v", response["error"])
	}
}

// Test Case 7: CreateAddress - Service Error
func TestCreateAddress_ServiceError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	address := models.Address{
		UserID:     1,
		Street:     "123 Main St",
		City:       "Cityville",
		State:      "Stateville",
		PostalCode: "12345",
	}

	mockService := &MockUserProfileService{
		SaveAddressFunc: func(addr models.Address) error {
			return errors.New("database error")
		},
	}

	handler := NewUserProfileHandler(mockService)

	body, _ := json.Marshal(address)
	c.Request = httptest.NewRequest("POST", "/users/address", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.CreateAddress(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Failed to save address" {
		t.Errorf("expected error message, got %v", response["error"])
	}
}

// Test Case 8: GetAddresses - Success
func TestGetAddresses_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedAddresses := []models.Address{
		{
			ID:         1,
			UserID:     2,
			Street:     "123 Main St",
			City:       "Cityville",
			State:      "Stateville",
			PostalCode: "12345",
		},
		{
			ID:         2,
			UserID:     2,
			Street:     "456 Side St",
			City:       "Townsville",
			State:      "Regionville",
			PostalCode: "67890",
		},
	}

	mockService := &MockUserProfileService{
		GetUserAdressesFunc: func(userId int) ([]models.Address, error) {
			if userId == 2 {
				return expectedAddresses, nil
			}
			return []models.Address{}, nil
		},
	}

	handler := NewUserProfileHandler(mockService)
	c.Request = httptest.NewRequest("GET", "/users/address", nil)

	// Act
	handler.GetAddresses(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var response []models.Address
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(response))
	}
}

// Test Case 9: GetAddresses - No Addresses
func TestGetAddresses_Empty(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{
		GetUserAdressesFunc: func(userId int) ([]models.Address, error) {
			return []models.Address{}, nil
		},
	}

	handler := NewUserProfileHandler(mockService)
	c.Request = httptest.NewRequest("GET", "/users/address", nil)

	// Act
	handler.GetAddresses(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var response []models.Address
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 0 {
		t.Errorf("expected 0 addresses, got %d", len(response))
	}
}

// Test Case 10: GetAddresses - Service Error
func TestGetAddresses_ServiceError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := &MockUserProfileService{
		GetUserAdressesFunc: func(userId int) ([]models.Address, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewUserProfileHandler(mockService)
	c.Request = httptest.NewRequest("GET", "/users/address", nil)

	// Act
	handler.GetAddresses(c)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !contains(response["error"].(string), "error while getting user adresses") {
		t.Errorf("expected error message containing 'error while getting user adresses', got %v", response["error"])
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
