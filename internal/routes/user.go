package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
)

func RegisterUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	followHandler *handlers.FollowHandler,
	authenticationMiddleware gin.HandlerFunc,
	rateLimiter *middleware.RateLimiter,
) {

	userRoutes := router.Group("/users")

	userRoutes.Use(
		authenticationMiddleware,
		rateLimiter.Middleware(defaultRateLimit),
	)

	{
		userRoutes.GET(
			"/me",
			authentication.RequireActiveUser(),
			userHandler.GetProfile,
		)

		userRoutes.PATCH(
			"/update",
			authentication.RequireActiveUser(),
			userHandler.UpdateProfile,
		)

		userRoutes.POST(
			"/activate",
			authentication.RequireUser(),
			userHandler.ActivateUser,
		)

		userRoutes.POST(
			"/deactivate",
			authentication.RequireActiveUser(),
			userHandler.DeactivateUser,
		)

		userRoutes.POST(
			"/:id/follow",
			authentication.RequireActiveUser(),
			followHandler.Follow,
		)

		userRoutes.DELETE(
			"/:id/unfollow",
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

		userRoutes.POST(
			"/profile-image",
			authentication.RequireActiveUser(),
			userHandler.CreateProfilePictureUpload,
		)

		userRoutes.POST(
			"/profile-image/confirm",
			authentication.RequireActiveUser(),
			userHandler.ConfirmProfilePictureUpload,
		)

		userRoutes.DELETE(
			"/profile-image",
			authentication.RequireActiveUser(),
			userHandler.DeleteProfilePicture,
		)
	}
}
