package handlers

import (
	"net/http"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (handler *MessageHandler) SendMessage(context *gin.Context) {
	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	var request dto.MessageCreate

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := handler.messageService.CreateMessage(user.ID, &request)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)

}

func (handler *MessageHandler) GetMessage(context *gin.Context) {
	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	messageID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message ID",
		})
		return
	}

	response, err := handler.messageService.GetMessage(messageID, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *MessageHandler) GetMessagesInChat(context *gin.Context) {
	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	chatID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat ID",
		})
		return
	}

	responses, err := handler.messageService.GetMessagesInChat(chatID, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, responses)
}
