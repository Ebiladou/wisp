package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/storage"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeactivateUser(userID uuid.UUID) error
	ActivateUser(userID uuid.UUID) error
	SearchUsers(userID uuid.UUID, query string) ([]*dto.PublicUserResponse, error)
	CreateProfilePictureUpload(ctx context.Context, userID uuid.UUID) (*storage.UploadResult, error)
	ConfirmProfilePictureUpload(ctx context.Context, userID uuid.UUID, imageID string) error
	DeleteProfilePicture(ctx context.Context, userID uuid.UUID) error
}

type DefaultUserService struct {
	userRepository repositories.UserRepository
	imageStorage   storage.ImageStorage
	logger         *slog.Logger
}

func NewUserService(
	userRepository repositories.UserRepository,
	imageStorage storage.ImageStorage,
	logger *slog.Logger,
) UserService {

	return &DefaultUserService{
		userRepository: userRepository,
		imageStorage:   imageStorage,
		logger:         logger,
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

func (service *DefaultUserService) CreateProfilePictureUpload(ctx context.Context, userID uuid.UUID) (*storage.UploadResult, error) {
	result, err := service.imageStorage.CreateUploadURL(ctx, userID.String())

	if err != nil {
		service.logger.Error(
			"failed to create profile picture upload URL",
			"user_id", userID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"create profile picture upload: %w",
			err,
		)
	}

	return result, nil
}

func (service *DefaultUserService) ConfirmProfilePictureUpload(ctx context.Context, userID uuid.UUID, imageID string) error {
	if imageID == "" {
		return errors.New("image ID is required")
	}

	// we get the image from cloudflare and verify it belongs to the user making the upload request
	image, err := service.imageStorage.GetImage(ctx, imageID)
	if err != nil {
		service.logger.Error(
			"failed to verify profile picture",
			"user_id", userID,
			"error", err,
		)

		return fmt.Errorf(
			"get profile picture: %w",
			err,
		)
	}

	if image.Creator != userID.String() {
		service.logger.Warn(
			"profile picture ownership validation failed",
		)
		return errors.New("image does not belong to user")
	}

	if !image.Uploaded {
		return errors.New("image upload is not complete")
	}

	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return fmt.Errorf(
			"find user: %w",
			err,
		)
	}

	oldImageID := user.ProfilePicture
	user.ProfilePicture = imageID

	if err := service.userRepository.UpdateProfilePicture(user); err != nil {
		service.logger.Error(
			"failed to update profile picture",
			"user_id", userID,
			"error", err,
		)
		return fmt.Errorf(
			"update profile picture: %w",
			err,
		)
	}

	if oldImageID != "" && oldImageID != imageID {
		if err := service.imageStorage.DeleteImage(
			ctx,
			oldImageID,
		); err != nil {
			service.logger.Error(
				"failed to delete old profile picture",
				"user_id", userID,
				"error", err,
			)
			return fmt.Errorf(
				"delete old profile picture: %w",
				err,
			)
		}
	}

	return nil
}

func (service *DefaultUserService) DeleteProfilePicture(ctx context.Context, userID uuid.UUID) error {
	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return fmt.Errorf(
			"find user: %w",
			err,
		)
	}

	if user.ProfilePicture == "" {
		return nil
	}

	imageID := user.ProfilePicture

	user.ProfilePicture = ""

	if err := service.userRepository.UpdateProfilePicture(user); err != nil {
		service.logger.Error(
			"failed to remove profile picture",
			"user_id", userID,
			"error", err,
		)

		return fmt.Errorf(
			"remove profile picture: %w",
			err,
		)
	}

	if err := service.imageStorage.DeleteImage(
		ctx,
		imageID,
	); err != nil {
		service.logger.Error(
			"failed to delete profile picture",
			"user_id", userID,
			"error", err,
		)

		return fmt.Errorf(
			"delete profile picture: %w",
			err,
		)
	}

	return nil
}
