package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Ebiladou/wisp/internal/models"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeactivateUser(user *models.User) error
	ActivateUser(user *models.User) error
}

type PostgreSQLUserRepository struct {
	database *gorm.DB
}

func NewPostgreSQLUserRepository(
	database *gorm.DB,
) UserRepository {

	return &PostgreSQLUserRepository{
		database: database,
	}
}

func (repository *PostgreSQLUserRepository) FindByID(id uuid.UUID) (*models.User, error) {

	var user models.User

	err := repository.database.
		First(&user, "id = ?", id).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repository *PostgreSQLUserRepository) FindByUsername(username string) (*models.User, error) {

	var user models.User

	err := repository.database.
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repository *PostgreSQLUserRepository) UpdateUser(user *models.User) error {

	return repository.database.Save(user).Error
}

func (repository *PostgreSQLUserRepository) DeactivateUser(user *models.User) error {

	user.Active = false
	user.DeletionRequested = true

	return repository.database.Save(user).Error
}

func (repository *PostgreSQLUserRepository) ActivateUser(user *models.User) error {

	user.Active = true
	user.DeletionRequested = false

	return repository.database.Save(user).Error
}
