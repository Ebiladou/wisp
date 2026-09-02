package routes

import (
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	rateLimiter *middleware.RateLimiter,
) {
	authRoutes := router.Group("/auth")

	authRoutes.POST(
		"/signup",
		rateLimiter.Middleware(authRateLimit),
		authHandler.CreateUser,
	)

	authRoutes.POST(
		"/confirm-email",
		rateLimiter.Middleware(authRateLimit),
		authHandler.ConfirmEmail,
	)

	authRoutes.POST(
		"/resend-confirmation",
		rateLimiter.Middleware(authRateLimit),
		authHandler.ResendConfirmation,
	)

	authRoutes.POST(
		"/login",
		rateLimiter.Middleware(authRateLimit),
		authHandler.Login,
	)

	authRoutes.POST(
		"/logout",
		rateLimiter.Middleware(authRateLimit),
		authHandler.Logout,
	)

	authRoutes.POST(
		"/forgot-password",
		rateLimiter.Middleware(sensitiveAuthLimit),
		authHandler.ForgotPassword,
	)

	authRoutes.POST(
		"/reset-password",
		rateLimiter.Middleware(sensitiveAuthLimit),
		authHandler.ResetPassword,
	)

	authRoutes.GET(
		"/users/:id",
		rateLimiter.Middleware(defaultRateLimit),
		authHandler.GetUserByID,
	)
}
