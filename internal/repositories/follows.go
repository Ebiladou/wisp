package repositories

import (
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowRepository interface {
	CreateFollow(follow *models.Follow) error
	RemoveFollow(followerID uuid.UUID, followingID uuid.UUID) error
	GetFollowing(userID uuid.UUID) ([]*models.User, error)
	GetFollowers(userID uuid.UUID) ([]*models.User, error)
}

type PostgreSQLFollowRepository struct {
	database *gorm.DB
}

func NewPostgreSQLFollowRepository(database *gorm.DB) FollowRepository {
	return &PostgreSQLFollowRepository{
		database: database,
	}
}

func (repository *PostgreSQLFollowRepository) CreateFollow(follow *models.Follow) error {
	return repository.database.Create(follow).Error
}

func (repository *PostgreSQLFollowRepository) RemoveFollow(followerID uuid.UUID, followingID uuid.UUID) error {
	return repository.database.
		Where(
			"follower_id = ? AND following_id = ?",
			followerID,
			followingID,
		).
		Delete(&models.Follow{}).Error
}

func (repository *PostgreSQLFollowRepository) GetFollowing(userID uuid.UUID) ([]*models.User, error) {

	var users []*models.User

	err := repository.database.
		Table("users").
		Joins(
			"JOIN follows ON follows.following_id = users.id",
		).
		Where(
			"follows.follower_id = ?",
			userID,
		).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *PostgreSQLFollowRepository) GetFollowers(userID uuid.UUID) ([]*models.User, error) {

	var users []*models.User

	err := repository.database.
		Table("users").
		Joins(
			"JOIN follows ON follows.follower_id = users.id",
		).
		Where(
			"follows.following_id = ?",
			userID,
		).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
