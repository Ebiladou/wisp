package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(
	userService services.UserService,
) *UserHandler {

	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) GetProfile(context *gin.Context) {

	user := context.MustGet(
		authentication.AuthenticatedUserKey,
	).(*models.User)

	response, err := handler.userService.GetProfile(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *UserHandler) UpdateProfile(context *gin.Context) {

	user := context.MustGet(
		authentication.AuthenticatedUserKey,
	).(*models.User)

	var request dto.UpdateUserRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := handler.userService.UpdateProfile(
		user.ID,
		request,
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *UserHandler) DeactivateUser(context *gin.Context) {
	user := context.MustGet(
		authentication.AuthenticatedUserKey,
	).(*models.User)

	err := handler.userService.DeactivateUser(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "account deactivated successfully",
	})
}

func (handler *UserHandler) ActivateUser(context *gin.Context) {

	user := context.MustGet(
		authentication.AuthenticatedUserKey,
	).(*models.User)

	err := handler.userService.ActivateUser(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "account activated successfully",
	})
}
