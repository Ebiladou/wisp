package handlers

import (
	"net/http"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) CreateUser(context *gin.Context) {

	var request dto.CreateUser

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := handler.authService.Create(request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (handler *AuthHandler) GetUserByID(context *gin.Context) {

	id := context.Param("id")

	userID, err := uuid.Parse(id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	response, err := handler.authService.GetByID(userID)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (handler *AuthHandler) ConfirmEmail(context *gin.Context) {

	var request dto.ConfirmEmailRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err = handler.authService.ConfirmEmail(request.Token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "email confirmed successfully",
	})
}

func (handler *AuthHandler) ResendConfirmation(context *gin.Context) {

	var request dto.ResendConfirmationRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err = handler.authService.ResendConfirmation(request.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "confirmation email sent successfully",
	})
}

func (handler *AuthHandler) ForgotPassword(context *gin.Context) {

	var request dto.ForgotPasswordRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err = handler.authService.ForgotPassword(request.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "if an account exists for that email, password reset link have been sent.",
	})
}

func (handler *AuthHandler) ResetPassword(context *gin.Context) {

	var request dto.ResetPasswordRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err = handler.authService.ResetPassword(
		request.Token,
		request.NewPassword,
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "password reset successfully",
	})
}
