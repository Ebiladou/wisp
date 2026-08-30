package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/services"
)

type BlockHandler struct {
	blockService services.BlockService
}

func NewBlockHandler(
	blockService services.BlockService,
) *BlockHandler {
	return &BlockHandler{
		blockService: blockService,
	}
}

func (handler *BlockHandler) BlockUser(context *gin.Context) {

	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	blockedID, err := uuid.Parse(
		context.Param("id"),
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	err = handler.blockService.BlockUser(user.ID, blockedID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "user blocked successfully",
	})
}

func (handler *BlockHandler) UnblockUser(context *gin.Context) {

	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	blockedID, err := uuid.Parse(
		context.Param("id"),
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	err = handler.blockService.UnblockUser(user.ID, blockedID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "user unblocked successfully",
	})
}

func (handler *BlockHandler) GetBlockedUsers(context *gin.Context) {

	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	response, err := handler.blockService.GetBlockedUsers(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get blocked users",
		})
		return
	}

	context.JSON(http.StatusOK, response)
}
