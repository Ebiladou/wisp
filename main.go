package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/database"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
	"github.com/Ebiladou/wisp/internal/observability"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/routes"
	"github.com/Ebiladou/wisp/internal/services"
	"github.com/Ebiladou/wisp/internal/storage"
)

func main() {

	logger := observability.NewLogger()

	logger.Info("Starting Wisp")

	applicationConfig, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.ConnectToDatabase(applicationConfig, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Repositories
	authRepository := repositories.NewPostgreSQLAuthRepository(db)
	tokenRepository := repositories.NewPostgreSQLTokenRepository(db)
	userRepository := repositories.NewPostgreSQLUserRepository(db)
	followRepository := repositories.NewPostgreSQLFollowRepository(db)
	blockRepository := repositories.NewPostgreSQLBlockRepository(db)
	chatRepository := repositories.NewPostgreSQLChatRepository(db)
	messageRepository := repositories.NewPostgreSQLMessageRepository(db)

	// storage

	imageStorage := storage.NewCloudflareImages(
		applicationConfig.CloudflareAccountID,
		applicationConfig.CloudflareAPIToken,
	)

	// Services

	accessPolicy := services.NewAccessPolicy(
		blockRepository,
	)

	authService := services.NewAuthService(
		authRepository,
		tokenRepository,
		applicationConfig,
		logger,
	)

	userService := services.NewUserService(
		userRepository,
		imageStorage,
		logger,
	)

	followService := services.NewFollowService(
		followRepository,
		userRepository,
		accessPolicy,
	)

	blockService := services.NewBlockService(
		blockRepository,
		userRepository,
		followRepository,
	)

	chatService := services.NewChatService(
		chatRepository,
	)

	messageService := services.NewMessageService(
		chatRepository,
		messageRepository,
		accessPolicy,
	)

	// Handlers
	authHandler := handlers.NewAuthHandler(
		authService,
		applicationConfig,
	)

	userHandler := handlers.NewUserHandler(
		userService,
	)

	followHandler := handlers.NewFollowHandler(
		followService,
	)

	blockHandler := handlers.NewBlockHandler(
		blockService,
	)

	chatHandler := handlers.NewChatHandler(
		chatService,
	)

	messageHandler := handlers.NewMessageHandler(
		messageService,
	)

	// Middlewares
	authenticationMiddleware := middleware.AuthMiddleware(
		authRepository,
		applicationConfig,
	)

	rateLimiter, err := middleware.NewRateLimiter("localhost:6379")
	if err != nil {
		logger.Error("failed to initialize rate limiter", "error", err)
		os.Exit(1)
	}
	defer rateLimiter.Close()

	// Routes
	router := gin.Default()
	router.Use(middleware.LoggingMiddleware(logger))

	routes.RegisterAuthRoutes(
		router,
		authHandler,
		rateLimiter,
	)

	routes.RegisterUserRoutes(
		router,
		userHandler,
		followHandler,
		authenticationMiddleware,
		rateLimiter,
	)

	routes.RegisterBlockRoutes(
		router,
		blockHandler,
		authenticationMiddleware,
		rateLimiter,
	)

	routes.RegisterChatRoutes(
		router,
		chatHandler,
		messageHandler,
		authenticationMiddleware,
		rateLimiter,
	)

	routes.RegisterMessageRoutes(
		router,
		messageHandler,
		authenticationMiddleware,
		rateLimiter,
	)

	err = router.Run(":" + applicationConfig.Port)
	if err != nil {
		logger.Error("failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}
