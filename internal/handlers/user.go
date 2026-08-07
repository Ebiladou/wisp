package handlers

import (
	"net/http"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) CreateUser(context *gin.Context) {

	var request dto.CreateUser

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := handler.userService.Create(request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (handler *UserHandler) GetUserByID(context *gin.Context) {

	id := context.Param("id")

	userID, err := uuid.Parse(id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	response, err := handler.userService.GetByID(userID)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}
