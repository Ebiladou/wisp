package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	ConfirmUser(user *models.User) error
	UpdatePassword(user *models.User) error

	// GetUser(user *models.User) (*models.User, error)
	// UpdateUser(user *models.User) (*models.User, error)
	// DeleteUser(user *models.User) error
}

type PostgreSQLAuthRepository struct {
	database *gorm.DB
}

func NewPostgreSQLAuthRepository(database *gorm.DB) AuthRepository {
	return &PostgreSQLAuthRepository{
		database: database,
	}
}

func (repository *PostgreSQLAuthRepository) CreateUser(user *models.User) error {
	return repository.database.Create(user).Error
}

func (repository *PostgreSQLAuthRepository) FindByID(id uuid.UUID) (*models.User, error) {

	var user models.User

	err := repository.database.
		Where("id = ?", id).
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

func (repository *PostgreSQLAuthRepository) FindByEmail(email string) (*models.User, error) {

	var user models.User

	err := repository.database.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repository *PostgreSQLAuthRepository) FindByUsername(username string) (*models.User, error) {

	var user models.User

	err := repository.database.Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil

}

func (repository *PostgreSQLAuthRepository) ConfirmUser(user *models.User) error {
	return repository.database.Model(user).Update("active", true).Error
}

func (repository *PostgreSQLAuthRepository) UpdatePassword(user *models.User) error {
	return repository.database.Model(user).Update("password", user.Password).Error
}
