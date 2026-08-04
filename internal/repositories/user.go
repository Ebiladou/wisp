package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Ebiladou/wisp/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
}

type PostgreSQLUserRepository struct {
	database *gorm.DB
}

func NewPostgreSQLUserRepository(database *gorm.DB) UserRepository {
	return &PostgreSQLUserRepository{
		database: database,
	}
}

func (repository *PostgreSQLUserRepository) Create(user *models.User) error {
	return repository.database.Create(user).Error
}

func (repository *PostgreSQLUserRepository) FindByID(id uuid.UUID) (*models.User, error) {

	var user models.User

	err := repository.database.First(&user, "id = ?", id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
