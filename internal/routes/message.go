package routes

import (
	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(
	router *gin.Engine,
	messageHandler *handlers.MessageHandler,
	authenticationMiddleware gin.HandlerFunc,
	rateLimiter *middleware.RateLimiter,
) {
	messageRoutes := router.Group("/messages")

	messageRoutes.Use(
		authenticationMiddleware,
		rateLimiter.Middleware(defaultRateLimit),
	)

	{
		messageRoutes.POST(
			"",
			authentication.RequireActiveUser(),
			messageHandler.SendMessage,
		)

		messageRoutes.GET(
			"/:id",
			authentication.RequireActiveUser(),
			messageHandler.GetMessage,
		)
	}

}
