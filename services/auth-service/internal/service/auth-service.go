package service

import (
	"errors"
	"log"

	"github.com/go-ecommerce-application/services/auth-service/internal/authentication"
	"github.com/go-ecommerce-application/services/auth-service/internal/dto"
	"github.com/go-ecommerce-application/services/auth-service/internal/models"
	"github.com/go-ecommerce-application/services/auth-service/internal/repository"
	"github.com/go-ecommerce-application/services/auth-service/internal/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(authData models.AuthUser) error
	Login(username, password string) (dto.LoginResponse, error)
	RefreshToken(refreshToken string) (newAccessToken string, err error)
	Logout(userID string) error
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (s *authService) Signup(authData models.AuthUser) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	authData.Id = uuid.String()
	hashPassword, err := utils.HashPassword(authData.Password)
	if err != nil {
		return err
	}
	authData.Password = hashPassword
	authData.Status = "active"
	err = s.authRepository.CreateUser(authData)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(email, password string) (dto.LoginResponse, error) {
	// Implement login logic here
	var loginResponse dto.LoginResponse
	var user dto.User
	if email == "" || password == "" {
		return dto.LoginResponse{}, errors.New("Email or Password can not be empty")
	}
	userData, err := s.authRepository.GetUserByEmail(email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if userData.Status != "active" {
		return dto.LoginResponse{}, errors.New("User is not active")
	}
	user.Id = userData.Id
	user.Email = userData.Email
	user.Role = userData.Role
	user.Status = userData.Status
	loginResponse.User = user

	if !utils.CheckPasswordHash(password, userData.Password) {
		return dto.LoginResponse{}, errors.New("email or password is incorrect")
	}
	accessToken, refreshToken, _, _, err := authentication.GenerateTokens(userData.Id, userData.Role)
	if err != nil {
		log.Println("Error while generating login token :", err)
		return dto.LoginResponse{}, err
	}

	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = refreshToken
	loginResponse.TokenType = "JWT"
	// loginResponse.ExpiresIn = accessExpiry
	// loginResponse.RefreshExpiresIn = refreshExpiry

	return loginResponse, nil
}

func (s *authService) RefreshToken(refreshToken string) (newAccessToken string, err error) {
	// Implement token refresh logic here
	return "", nil
}

func (s *authService) Logout(userID string) error {
	// Implement logout logic here
	return nil
}
