package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/services"
)

type FollowHandler struct {
	followService services.FollowService
}

func NewFollowHandler(
	followService services.FollowService,
) *FollowHandler {

	return &FollowHandler{
		followService: followService,
	}
}

func (handler *FollowHandler) Follow(context *gin.Context) {

	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	followingID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	err = handler.followService.FollowUser(user.ID, followingID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "user followed successfully"})
}

func (handler *FollowHandler) Unfollow(context *gin.Context) {

	user := context.MustGet(authentication.AuthenticatedUserKey).(*models.User)

	followingID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	err = handler.followService.UnfollowUser(user.ID, followingID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user unfollowed successfully"})
}

func (handler *FollowHandler) GetFollowing(context *gin.Context) {

	userID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	response, err := handler.followService.GetFollowing(userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *FollowHandler) GetFollowers(context *gin.Context) {

	userID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	response, err := handler.followService.GetFollowers(userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}
