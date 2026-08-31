package services

import (
	"errors"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/google/uuid"
)

type BlockService interface {
	BlockUser(blockerID uuid.UUID, blockedID uuid.UUID) error
	UnblockUser(blockerID uuid.UUID, blockedID uuid.UUID) error
	GetBlockedUsers(blockerID uuid.UUID) ([]*dto.PublicUserResponse, error)
}

type DefaultBlockService struct {
	blockRepository  repositories.BlockRepository
	userRepository   repositories.UserRepository
	followRepository repositories.FollowRepository
}

func NewBlockService(
	blockRepository repositories.BlockRepository,
	userRepository repositories.UserRepository,
	followRepository repositories.FollowRepository,
) BlockService {

	return &DefaultBlockService{
		blockRepository:  blockRepository,
		userRepository:   userRepository,
		followRepository: followRepository,
	}
}

func (service *DefaultBlockService) BlockUser(blockerID uuid.UUID, blockedID uuid.UUID) error {

	user, err := service.userRepository.FindByID(blockedID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	exists, err := service.blockRepository.BlockExists(
		blockerID,
		blockedID,
	)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user is already blocked")
	}

	block := models.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	err = service.blockRepository.CreateBlock(&block)

	if err != nil {
		return err
	}

	err = service.followRepository.RemoveFollow(blockerID, blockedID)

	if err != nil {
		return err
	}

	err = service.followRepository.RemoveFollow(blockedID, blockerID)

	if err != nil {
		return err
	}

	return nil
}

func (service *DefaultBlockService) UnblockUser(blockerID uuid.UUID, blockedID uuid.UUID) error {

	return service.blockRepository.RemoveBlock(blockerID, blockedID)
}

func (service *DefaultBlockService) GetBlockedUsers(blockerID uuid.UUID) ([]*dto.PublicUserResponse, error) {

	users, err := service.blockRepository.GetBlockedUsers(blockerID)

	if err != nil {
		return nil, err
	}

	var responses []*dto.PublicUserResponse

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
