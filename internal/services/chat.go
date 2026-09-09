package services

import (
	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/google/uuid"
)

type ChatService interface {
	GetUserChat(userID uuid.UUID, chatID uuid.UUID) (*dto.ChatResponse, error)
	GetUserChats(userID uuid.UUID) ([]*dto.ChatResponse, error)
}

type DefaultChatService struct {
	chatRepository repositories.ChatRepository
}

func NewChatService(
	chatRepository repositories.ChatRepository,
) ChatService {
	return &DefaultChatService{
		chatRepository: chatRepository,
	}
}

func (service *DefaultChatService) GetUserChat(userID uuid.UUID, chatID uuid.UUID) (*dto.ChatResponse, error) {
	chat, err := service.chatRepository.GetChat(userID, chatID)

	if err != nil {
		return nil, err
	}

	return &dto.ChatResponse{
		ID:        chat.ID.String(),
		UserOneID: chat.UserOneID.String(),
		UserTwoID: chat.UserTwoID.String(),
		CreatedAt: chat.CreatedAt,
	}, nil
}

func (service *DefaultChatService) GetUserChats(userID uuid.UUID) ([]*dto.ChatResponse, error) {
	chats, err := service.chatRepository.GetChats(userID)

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ChatResponse, 0, len(chats))

	for _, chat := range chats {
		responses = append(responses, &dto.ChatResponse{
			ID:        chat.ID.String(),
			UserOneID: chat.UserOneID.String(),
			UserTwoID: chat.UserTwoID.String(),
			CreatedAt: chat.CreatedAt,
		})
	}

	return responses, nil
}
