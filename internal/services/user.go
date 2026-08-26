package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/repositories"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type DefaultUserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {

	return &DefaultUserService{
		userRepository: userRepository,
	}
}

func (service *DefaultUserService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {

	user, err := service.userRepository.FindByID(userID)

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

func (service *DefaultUserService) UpdateProfile(userID uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error) {

	user, err := service.userRepository.FindByID(userID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Username != "" {

		existingUser, err :=
			service.userRepository.FindByUsername(request.Username)

		if err != nil {
			return nil, err
		}

		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("username already taken")
		}

		user.Username = request.Username
	}

	if request.ProfilePicture != "" {
		user.ProfilePicture = request.ProfilePicture
	}

	err = service.userRepository.UpdateUser(user)

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
