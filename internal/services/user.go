package services

import (
	"errors"
	"log"
	"time"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/utils"
)

type UserService interface {
	Create(request dto.CreateUser) (*dto.UserResponse, error)
	// GetByID(id uuid.UUID) (*dto.UserResponse, error)
}

type DefaultUserService struct {
	userRepository  repositories.UserRepository
	tokenRepository repositories.TokenRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
	tokenRepository repositories.TokenRepository,
) UserService {
	return &DefaultUserService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (service *DefaultUserService) Create(request dto.CreateUser) (*dto.UserResponse, error) {

	existingUser, err := service.userRepository.FindByEmail(request.Email)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = service.userRepository.FindByUsername(request.Username)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    hashedPassword,
		DateOfBirth: request.DateOfBirth,
		Active:      false,
	}

	err = service.userRepository.CreateUser(&user)

	if err != nil {
		return nil, err
	}

	rawToken, err := utils.GenerateToken()

	if err != nil {
		return nil, err
	}

	hashedToken := utils.HashToken(rawToken)

	token := models.Token{
		UserID:    user.ID,
		Token:     hashedToken,
		TokenType: models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = service.tokenRepository.Create(&token)

	if err != nil {
		return nil, err
	}

	// for testing, please. we'll move to an email service soon enough when my bag is up.
	log.Printf(
		"email verification token for user %s: %s",
		user.Email,
		rawToken,
	)

	response := &dto.UserResponse{
		ID:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Active:         user.Active,
	}

	return response, nil
}
