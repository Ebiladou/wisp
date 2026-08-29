package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
)

func RegisterUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	followHandler *handlers.FollowHandler,
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

		userRoutes.POST(
			"/deactivate",
			authentication.RequireActiveUser(),
			userHandler.DeactivateUser,
		)

		userRoutes.POST(
			"/activate",
			authentication.RequireUser(),
			userHandler.ActivateUser,
		)

		userRoutes.POST(
			"/:id/follow",
			authentication.RequireActiveUser(),
			followHandler.Follow,
		)

		userRoutes.DELETE(
			"/:id/follow",
			authentication.RequireActiveUser(),
			followHandler.Unfollow,
		)

		userRoutes.GET(
			"",
			authentication.RequireActiveUser(),
			userHandler.SearchUsers,
		)

		userRoutes.GET(
			"/:id/followers",
			authentication.RequireActiveUser(),
			followHandler.GetFollowers,
		)

		userRoutes.GET(
			"/:id/following",
			authentication.RequireActiveUser(),
			followHandler.GetFollowing,
		)
	}
}
