package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
)

type UserService interface {
	Create(request dto.CreateUser) (*dto.UserResponse, error)
	GetByID(id uuid.UUID) (*dto.UserResponse, error)
}

type DefaultUserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &DefaultUserService{
		userRepository: userRepository,
	}
}

func (service *DefaultUserService) Create(
	request dto.CreateUser,
) (*dto.UserResponse, error) {

	user := models.User{
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    request.Password,
		DateOfBirth: request.DateOfBirth,
		Active:      false,
	}

	err := service.userRepository.Create(&user)

	if err != nil {
		return nil, err
	}

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

func (service *DefaultUserService) GetByID(
	id uuid.UUID,
) (*dto.UserResponse, error) {

	user, err := service.userRepository.FindByID(id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

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
