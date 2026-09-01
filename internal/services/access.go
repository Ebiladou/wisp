package services

import (
	"errors"

	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/google/uuid"
)

type AccessPolicy interface {
	CanAccess(userID uuid.UUID, targetUserID uuid.UUID) error
}

type DefaultAccessPolicy struct {
	blockRepository repositories.BlockRepository
}

func NewAccessPolicy(
	blockRepository repositories.BlockRepository,
) AccessPolicy {
	return &DefaultAccessPolicy{
		blockRepository: blockRepository,
	}
}

func (policy *DefaultAccessPolicy) CanAccess(userID uuid.UUID, targetUserID uuid.UUID) error {

	blocked, err := policy.blockRepository.IsBlocked(userID, targetUserID)

	if err != nil {
		return err
	}

	if blocked {
		return errors.New("interaction not allowed between blocked users")
	}

	return nil
}
