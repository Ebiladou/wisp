package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
)

func RegisterChatRoutes(
	router *gin.Engine,
	chatHandler *handlers.ChatHandler,
	authenticationMiddleware gin.HandlerFunc,
	rateLimiter *middleware.RateLimiter,
) {

	chatRoutes := router.Group("/chats")

	chatRoutes.Use(
		authenticationMiddleware,
		rateLimiter.Middleware(defaultRateLimit),
	)

	{
		chatRoutes.GET(
			"",
			authentication.RequireActiveUser(),
			chatHandler.GetChats,
		)

		chatRoutes.GET(
			"/:id",
			authentication.RequireActiveUser(),
			chatHandler.GetChat,
		)

	}
}
