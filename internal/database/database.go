package database

import (
	"fmt"
	"log"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(applicationConfig *config.Config) (*gorm.DB, error) {

	db, err := gorm.Open(
		postgres.Open(applicationConfig.DatabaseURL),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Successfully connected to database")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migrated successfully.")

	return db, nil
}
