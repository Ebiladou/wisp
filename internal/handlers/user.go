package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	value, exists := context.Get("user")

	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user is required",
		})
		return
	}

	user, ok := value.(*models.User)

	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid authenticated user",
		})
		return
	}

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

	value, exists := context.Get("user")

	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user is required",
		})
		return
	}

	user, ok := value.(*models.User)

	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid authenticated user",
		})
		return
	}

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
