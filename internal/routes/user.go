package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
)

func RegisterUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	authenticationMiddleware gin.HandlerFunc,
) {

	userRoutes := router.Group("/users")

	userRoutes.Use(authenticationMiddleware)

	{
		userRoutes.GET(
			"/me",
			authentication.RequireActiveUser(),
			userHandler.GetProfile,
		)

		userRoutes.PATCH(
			"/me",
			authentication.RequireActiveUser(),
			userHandler.UpdateProfile,
		)
	}
}
