package handlers

import (
	"net/http"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService services.ChatService
}

func NewChatHandler(
	chatService services.ChatService,
) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (handler *ChatHandler) GetChat(context *gin.Context) {
	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	chatID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
		return
	}

	response, err := handler.chatService.GetUserChat(user.ID, chatID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *ChatHandler) GetChats(context *gin.Context) {
	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	response, err := handler.chatService.GetUserChats(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}
