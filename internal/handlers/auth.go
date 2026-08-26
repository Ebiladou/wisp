package handlers

import (
	"net/http"
	"time"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/services"
	"github.com/Ebiladou/wisp/internal/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	config      *config.Config
}

func NewAuthHandler(
	authService services.AuthService,
	applicationConfig *config.Config,
) *AuthHandler {

	return &AuthHandler{
		authService: authService,
		config:      applicationConfig,
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

func (handler *AuthHandler) Login(context *gin.Context) {

	var request dto.LoginRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	accessToken, refreshToken, err := handler.authService.Login(request)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.SetCookie(
		"access_token",
		accessToken,
		int((time.Duration(
			handler.config.JwtExpiryMins,
		) * time.Minute).Seconds()),
		"/",
		"",
		true,
		true,
	)

	context.SetCookie(
		"refresh_token",
		refreshToken,
		int(utils.RefreshTokenExpiry.Seconds()),
		"/",
		"",
		true,
		true,
	)

	context.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}

func (handler *AuthHandler) Logout(context *gin.Context) {

	refreshToken, _ := context.Cookie("refresh_token")

	err := handler.authService.Logout(refreshToken)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	context.SetCookie("access_token", "", -1, "/", "", true, true)

	context.SetCookie("refresh_token", "", -1, "/", "", true, true)

	context.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
