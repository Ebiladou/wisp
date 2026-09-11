package repositories

import (
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// i don't think deleting chat is necessary in the event of blocking. when a user is blocked, they simply are no longer permitted to have further chat sessions with the bloker. once unblocked, messages in chats can resume as normal, but follow relation is terminated as originally designed. will look into this, fingers crossed it's accurate.

type ChatRepository interface {
	CreateChat(chat *models.Chat) error
	GetChat(userID uuid.UUID, chatID uuid.UUID) (*models.Chat, error)
	GetChats(userID uuid.UUID) ([]*models.Chat, error)
	ChatBetweenUsers(userOneID uuid.UUID, userTwoID uuid.UUID) (*models.Chat, error)
}

type PostgreSQLChatRepository struct {
	database *gorm.DB
}

func NewPostgreSQLChatRepository(database *gorm.DB) ChatRepository {
	return &PostgreSQLChatRepository{
		database: database,
	}
}

func (repository *PostgreSQLChatRepository) CreateChat(chat *models.Chat) error {
	return repository.database.Create(chat).Error
}

func (repository *PostgreSQLChatRepository) GetChat(userID uuid.UUID, chatID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat

	err := repository.database.Where(
		"id = ? AND (user_one_id = ? OR user_two_id = ?)",
		chatID, userID, userID,
	).First(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil

}

func (repository *PostgreSQLChatRepository) GetChats(userID uuid.UUID) ([]*models.Chat, error) {
	var chats []*models.Chat

	err := repository.database.Where(
		"user_one_id = ? OR user_two_id = ?",
		userID,
		userID,
	).Find(&chats).Error

	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (repository *PostgreSQLChatRepository) ChatBetweenUsers(userOneID uuid.UUID, userTwoID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat

	err := repository.database.Where(
		"(user_one_id = ? AND user_two_id = ?) OR "+
			"(user_two_id = ? AND user_one_id = ?)",
		userOneID,
		userTwoID,
		userTwoID,
		userOneID,
	).First(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}
