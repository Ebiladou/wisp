package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/database"
	"github.com/Ebiladou/wisp/internal/handlers"
	"github.com/Ebiladou/wisp/internal/middleware"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/routes"
	"github.com/Ebiladou/wisp/internal/services"
)

func main() {

	applicationConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectToDatabase(applicationConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	authRepository := repositories.NewPostgreSQLAuthRepository(db)
	tokenRepository := repositories.NewPostgreSQLTokenRepository(db)
	userRepository := repositories.NewPostgreSQLUserRepository(db)
	followRepository := repositories.NewPostgreSQLFollowRepository(db)
	blockRepository := repositories.NewPostgreSQLBlockRepository(db)

	// Services

	accessPolicy := services.NewAccessPolicy(
		blockRepository,
	)

	authService := services.NewAuthService(
		authRepository,
		tokenRepository,
		applicationConfig,
	)

	userService := services.NewUserService(
		userRepository,
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

	// Middlewares
	authenticationMiddleware := middleware.AuthMiddleware(
		authRepository,
		applicationConfig,
	)

	rateLimiter, err := middleware.NewRateLimiter("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer rateLimiter.Close()

	// Routes
	router := gin.Default()

	routes.RegisterAuthRoutes(
		router,
		authHandler,
		rateLimiter,
	)

	routes.RegisterUserRoutes(
		router,
		userHandler,
		followHandler,
		blockHandler,
		authenticationMiddleware,
		rateLimiter,
	)

	err = router.Run(":" + applicationConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
}
