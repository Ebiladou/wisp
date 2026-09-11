package services

import (
	"errors"
	"time"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/google/uuid"
)

type MessageService interface {
	CreateMessage(userID uuid.UUID, message *dto.MessageCreate) (*dto.MessageResponse, error)
	GetMessage(messageID uuid.UUID, userID uuid.UUID) (*dto.MessageResponse, error)
	GetMessagesInChat(chatID uuid.UUID, userID uuid.UUID) ([]*dto.MessageResponse, error)
}

type DefaultMessageService struct {
	chatRepository    repositories.ChatRepository
	messageRepository repositories.MessageRepository
	accessPolicy      AccessPolicy
}

func NewMessageService(
	chatRepository repositories.ChatRepository,
	messageRepository repositories.MessageRepository,
	accessPolicy AccessPolicy,
) MessageService {
	return &DefaultMessageService{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
		accessPolicy:      accessPolicy,
	}
}

func (service *DefaultMessageService) CreateMessage(userID uuid.UUID, message *dto.MessageCreate) (*dto.MessageResponse, error) {

	if userID == message.RecipientID {
		return nil, errors.New("cannot text yourself")
	}

	if err := service.accessPolicy.CanAccess(userID, message.RecipientID); err != nil {
		return nil, err
	}

	const maximumExpiryMinutes = 30 * 24 * 60

	if message.ExpiryMinutes < 1 {
		return nil, errors.New("message expiry must be at least 1 minute")
	}

	if message.ExpiryMinutes > maximumExpiryMinutes {
		return nil, errors.New("message expiry must be within 30 days")
	}

	switch message.Type {
	case models.MessageTypeText:
		if message.Content == nil || *message.Content == "" {
			return nil, errors.New("cannot send empty text")
		}

		if message.MediaID != nil {
			return nil, errors.New("text message cannot contain media")
		}

	case models.MessageTypeImage,
		models.MessageTypeVideo,
		models.MessageTypeAudio:

		if message.MediaID == nil {
			return nil, errors.New("media message must contain media")
		}

		if message.Content != nil {
			return nil, errors.New("media message cannot contain text")
		}

	default:
		return nil, errors.New("invalid message type")
	}

	chat, err := service.chatRepository.ChatBetweenUsers(userID, message.RecipientID)

	if err != nil {
		return nil, err
	}

	if chat == nil {
		chat = &models.Chat{
			UserOneID: userID,
			UserTwoID: message.RecipientID,
		}

		if err := service.chatRepository.CreateChat(chat); err != nil {
			return nil, err
		}
	}

	expiresAt := time.Now().Add(time.Duration(message.ExpiryMinutes) * time.Minute)

	newMessage := models.Message{
		ChatID:    chat.ID,
		SenderID:  userID,
		Type:      message.Type,
		Content:   message.Content,
		MediaID:   message.MediaID,
		ExpiresAt: expiresAt,
	}

	if err := service.messageRepository.CreateMessage(&newMessage); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		ID:          newMessage.ID,
		ChatID:      newMessage.ChatID,
		SenderID:    newMessage.SenderID,
		RecipientID: message.RecipientID,
		Type:        newMessage.Type,
		Content:     newMessage.Content,
		MediaID:     newMessage.MediaID,
		CreatedAt:   newMessage.CreatedAt,
		ExpiresAt:   expiresAt,
	}, nil

}

func (service *DefaultMessageService) GetMessage(messageID uuid.UUID, userID uuid.UUID) (*dto.MessageResponse, error) {

	message, err := service.messageRepository.GetMessage(messageID, userID)

	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, err
	}

	chat, err := service.chatRepository.GetChat(userID, message.ChatID)

	if err != nil {
		return nil, err
	}

	var recipientID uuid.UUID

	if userID == chat.UserOneID {
		recipientID = chat.UserTwoID
	} else {
		recipientID = chat.UserOneID
	}

	return &dto.MessageResponse{
		ID:          message.ID,
		ChatID:      message.ChatID,
		SenderID:    message.SenderID,
		RecipientID: recipientID,
		Type:        message.Type,
		Content:     message.Content,
		MediaID:     message.MediaID,
		CreatedAt:   message.CreatedAt,
		ExpiresAt:   message.ExpiresAt,
	}, nil
}

func (service *DefaultMessageService) GetMessagesInChat(chatID uuid.UUID, userID uuid.UUID) ([]*dto.MessageResponse, error) {
	messages, err := service.messageRepository.GetChatMessages(chatID, userID)

	if err != nil {
		return nil, err
	}

	chat, err := service.chatRepository.GetChat(userID, chatID)

	if err != nil {
		return nil, err
	}

	var recipientID uuid.UUID

	if userID == chat.UserOneID {
		recipientID = chat.UserTwoID
	} else {
		recipientID = chat.UserOneID
	}

	responses := make([]*dto.MessageResponse, 0, len(messages))

	for _, message := range messages {
		responses = append(responses, &dto.MessageResponse{
			ID:          message.ID,
			ChatID:      message.ChatID,
			SenderID:    message.SenderID,
			RecipientID: recipientID,
			Type:        message.Type,
			Content:     message.Content,
			MediaID:     message.MediaID,
			CreatedAt:   message.CreatedAt,
			ExpiresAt:   message.ExpiresAt,
		})
	}

	return responses, nil
}
