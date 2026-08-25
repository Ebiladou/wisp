package routes

import (
	"github.com/Ebiladou/wisp/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
) {
	userRoutes := router.Group("/users")

	{
		userRoutes.POST("/register", userHandler.CreateUser)
		// userRoutes.GET("/:id", userHandler.GetUserByID)
	}
}
