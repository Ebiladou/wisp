package repositories

import (
	"time"

	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *models.Message) error
	GetMessage(messageID uuid.UUID, userID uuid.UUID) (*models.Message, error) // now i'm wondering if recipient id should be in the message model so we can easily check ownership between sender and recipient id without having to query chat. seems like a less expensive operation.
	GetChatMessages(chatID uuid.UUID, userID uuid.UUID) ([]*models.Message, error)
	GetExpiredMessages(expiryTime time.Time) ([]*models.Message, error)
}

type PostgreSQLMessageRepository struct {
	database *gorm.DB
}

func NewPostgreSQLMessageRepository(database *gorm.DB) MessageRepository {
	return &PostgreSQLMessageRepository{
		database: database,
	}
}

func (repository *PostgreSQLMessageRepository) CreateMessage(message *models.Message) error {
	return repository.database.Create(message).Error
}

func (repository *PostgreSQLMessageRepository) GetMessage(messageID uuid.UUID, userID uuid.UUID) (*models.Message, error) {

	var message models.Message

	err := repository.database.
		Table("messages").
		Joins(
			"JOIN chats ON chats.id = messages.chat_id",
		).
		Where(
			"messages.id = ? AND "+
				"(chats.user_one_id = ? OR chats.user_two_id = ?) AND "+
				"messages.expires_at > ? AND "+
				"messages.deleted = ?",
			messageID,
			userID,
			userID,
			time.Now(),
			false,
		).
		First(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (repository *PostgreSQLMessageRepository) GetChatMessages(chatID uuid.UUID, userID uuid.UUID) ([]*models.Message, error) {

	var messages []*models.Message

	err := repository.database.
		Table("messages").
		Joins(
			"JOIN chats ON chats.id = messages.chat_id",
		).
		Where(
			"chats.id = ? AND "+
				"(chats.user_one_id = ? OR chats.user_two_id = ?) AND "+
				"messages.expires_at > ? AND "+
				"messages.deleted = ?",
			chatID,
			userID,
			userID,
			time.Now(),
			false,
		).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (repository *PostgreSQLMessageRepository) GetExpiredMessages(expiryTime time.Time) ([]*models.Message, error) {

	var messages []*models.Message

	err := repository.database.
		Where(
			"expires_at <= ? AND deleted = ?",
			expiryTime,
			false,
		).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
