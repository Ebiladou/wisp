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
	DeactivateUser(userID uuid.UUID) error
	ActivateUser(userID uuid.UUID) error
	SearchUsers(userID uuid.UUID, query string) ([]*dto.PublicUserResponse, error)
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

func (service *DefaultUserService) DeactivateUser(userID uuid.UUID) error {

	user, err := service.userRepository.FindByID(userID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return service.userRepository.DeactivateUser(user)
}

func (service *DefaultUserService) ActivateUser(userID uuid.UUID) error {

	user, err := service.userRepository.FindByID(userID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return service.userRepository.ActivateUser(user)
}

func (service *DefaultUserService) SearchUsers(userID uuid.UUID, query string) ([]*dto.PublicUserResponse, error) {

	users, err := service.userRepository.SearchUsers(userID, query)

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.PublicUserResponse, 0)

	for _, user := range users {
		responses = append(responses, &dto.PublicUserResponse{
			ID:             user.ID.String(),
			Name:           user.Name,
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
		})
	}

	return responses, nil
}
