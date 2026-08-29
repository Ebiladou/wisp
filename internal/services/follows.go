package services

import (
	"errors"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/google/uuid"
)

type FollowService interface {
	FollowUser(followerID uuid.UUID, followingID uuid.UUID) error
	UnfollowUser(followerID uuid.UUID, followingID uuid.UUID) error
	GetFollowing(userID uuid.UUID) ([]*dto.UserResponse, error)
	GetFollowers(userID uuid.UUID) ([]*dto.UserResponse, error)
}

type DefaultFollowService struct {
	followRepository repositories.FollowRepository
	userRepository   repositories.UserRepository
}

func NewFollowService(
	followRepository repositories.FollowRepository,
	userRepository repositories.UserRepository,
) FollowService {

	return &DefaultFollowService{
		followRepository: followRepository,
		userRepository:   userRepository,
	}
}

func (service *DefaultFollowService) FollowUser(followerID uuid.UUID, followingID uuid.UUID) error {

	user, err := service.userRepository.FindByID(followingID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	follow := models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return service.followRepository.CreateFollow(&follow)
}

func (service *DefaultFollowService) UnfollowUser(followerID uuid.UUID, followingID uuid.UUID) error {

	return service.followRepository.RemoveFollow(followerID, followingID)
}

func (service *DefaultFollowService) GetFollowing(userID uuid.UUID) ([]*dto.UserResponse, error) {

	users, err := service.followRepository.GetFollowing(userID)

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		responses = append(responses, &dto.UserResponse{
			ID:             user.ID.String(),
			Name:           user.Name,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			Active:         user.Active,
		})
	}

	return responses, nil
}

func (service *DefaultFollowService) GetFollowers(userID uuid.UUID) ([]*dto.UserResponse, error) {

	users, err := service.followRepository.GetFollowers(userID)

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		responses = append(responses, &dto.UserResponse{
			ID:             user.ID.String(),
			Name:           user.Name,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			Active:         user.Active,
		})
	}

	return responses, nil
}
