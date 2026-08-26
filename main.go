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

	AuthRepository := repositories.NewPostgreSQLAuthRepository(db)
	tokenRepository := repositories.NewPostgreSQLTokenRepository(db)

	AuthService := services.NewAuthService(
		AuthRepository,
		tokenRepository,
	)

	AuthHandler := handlers.NewAuthHandler(AuthService)

	authenticationMiddleware :=
		middleware.NewAuthenticationMiddleware(
			applicationConfig,
			AuthRepository,
		)

	router := gin.Default()

	routes.RegisterUserRoutes(
		router,
		AuthHandler,
	)

	_ = authenticationMiddleware

	err = router.Run(":" + applicationConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
}
