package repositories

import (
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlockRepository interface {
	CreateBlock(block *models.Block) error
	RemoveBlock(blockerID uuid.UUID, blockedID uuid.UUID) error
	GetBlockedUsers(blockerID uuid.UUID) ([]*models.User, error)
}

type PostgreSQLBlockRepository struct {
	database *gorm.DB
}

func NewPostgreSQLBlockRepository(
	database *gorm.DB,
) BlockRepository {
	return &PostgreSQLBlockRepository{
		database: database,
	}
}

func (repository *PostgreSQLBlockRepository) CreateBlock(block *models.Block) error {
	return repository.database.Create(block).Error
}

func (repository *PostgreSQLBlockRepository) RemoveBlock(blockerID uuid.UUID, blockedID uuid.UUID) error {

	return repository.database.
		Where(
			"blocker_id = ? AND blocked_id = ?",
			blockerID,
			blockedID,
		).
		Delete(&models.Block{}).
		Error
}

func (repository *PostgreSQLBlockRepository) GetBlockedUsers(blockerID uuid.UUID) ([]*models.User, error) {

	var users []*models.User

	err := repository.database.
		Joins(
			"JOIN blocks ON blocks.blocked_id = users.id",
		).
		Where(
			"blocks.blocker_id = ?",
			blockerID,
		).
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
