package service

import (
	"errors"
	"log"

	"github.com/go-ecommerce-application/pkg/auth"
	"github.com/go-ecommerce-application/services/auth-service/internal/models"
	"github.com/go-ecommerce-application/services/auth-service/internal/repository"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(authData models.AuthUser) error
	Login(username, password string) (auth.LoginResponse, error)
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
	hashPassword, err := auth.HashPassword(authData.Password)
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

func (s *authService) Login(email, password string) (auth.LoginResponse, error) {
	// Implement login logic here
	var loginResponse auth.LoginResponse
	var user auth.User
	if email == "" || password == "" {
		return auth.LoginResponse{}, errors.New("Email or Password can not be empty")
	}
	userData, err := s.authRepository.GetUserByEmail(email)
	if err != nil {
		return auth.LoginResponse{}, err
	}
	if userData.Status != "active" {
		return auth.LoginResponse{}, errors.New("User is not active")
	}
	user.Id = userData.Id
	user.Email = userData.Email
	user.Role = userData.Role
	user.Status = userData.Status
	loginResponse.User = user

	if !auth.CheckPasswordHash(password, userData.Password) {
		return auth.LoginResponse{}, errors.New("email or password is incorrect")
	}
	accessToken, refreshToken, _, _, err := auth.GenerateTokens(userData.Id, userData.Role)
	if err != nil {
		log.Println("Error while generating login token :", err)
		return auth.LoginResponse{}, err
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
