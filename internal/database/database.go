package database

import (
	"fmt"
	"log/slog"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(applicationConfig *config.Config, logger *slog.Logger) (*gorm.DB, error) {

	db, err := gorm.Open(
		postgres.Open(applicationConfig.DatabaseURL),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Successfully connected to database")

	err = db.AutoMigrate(
		&models.User{},
		&models.Follow{},
		&models.Block{},
		&models.Token{},
		&models.Chat{},
		&models.Message{},
		&models.Media{},
		&models.ViewingSession{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	logger.Info("Database migrated successfully.")

	return db, nil
}
