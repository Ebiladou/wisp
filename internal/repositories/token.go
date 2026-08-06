package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Ebiladou/wisp/internal/models"
)

type TokenRepository interface {
	Create(token *models.Token) error

	FindActiveToken(
		userID uuid.UUID,
		tokenType models.TokenType,
	) (*models.Token, error)
}

type PostgreSQLTokenRepository struct {
	database *gorm.DB
}

func NewPostgreSQLTokenRepository(
	database *gorm.DB,
) TokenRepository {

	return &PostgreSQLTokenRepository{
		database: database,
	}
}

func (repository *PostgreSQLTokenRepository) Create(
	token *models.Token,
) error {

	return repository.database.Create(token).Error
}

func (repository *PostgreSQLTokenRepository) FindActiveToken(
	userID uuid.UUID,
	tokenType models.TokenType,
) (*models.Token, error) {

	var token models.Token

	err := repository.database.
		Where(
			"user_id = ? AND token_type = ? AND used_at IS NULL",
			userID,
			tokenType,
		).
		First(&token).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}
