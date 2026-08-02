package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/database"
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

	router := gin.Default()

	router.Run(":" + applicationConfig.Port)

	_ = db
}
