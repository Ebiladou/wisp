package routes

import (
	"github.com/Ebiladou/wisp/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
) {
	authRoutes := router.Group("/auth")

	{
		authRoutes.POST("/register", authHandler.CreateUser)
		authRoutes.GET("/users/:id", authHandler.GetUserByID)
		authRoutes.POST("/confirm-email", authHandler.ConfirmEmail)
		authRoutes.POST("/resend-confirmation", authHandler.ResendConfirmation)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
	}
}
