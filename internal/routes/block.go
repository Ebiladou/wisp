package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
)

func RegisterBlockRoutes(
	router *gin.Engine,
	blockHandler *handlers.BlockHandler,
	authenticationMiddleware gin.HandlerFunc,
	rateLimiter *middleware.RateLimiter,
) {

	blockRoutes := router.Group("/blocks")

	blockRoutes.Use(
		authenticationMiddleware,
		rateLimiter.Middleware(defaultRateLimit),
	)

	{
		blockRoutes.GET(
			"",
			authentication.RequireActiveUser(),
			blockHandler.GetBlockedUsers,
		)

		blockRoutes.POST(
			"/:id",
			authentication.RequireActiveUser(),
			blockHandler.BlockUser,
		)

		blockRoutes.DELETE(
			"/:id",
			authentication.RequireActiveUser(),
			blockHandler.UnblockUser,
		)

	}
}
