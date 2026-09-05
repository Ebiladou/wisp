package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DbUser              string
	DbPassword          string
	DbName              string
	DatabaseURL         string
	Host                string
	JwtSecret           string
	JwtExpiryMins       int
	Algorithm           string
	CloudflareAccountID string
	CloudflareAPIToken  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load env file")
	}
	log.Printf("config loaded")

	jwtExpiryMins, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINS"))
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT_EXPIRY_MINS")
	}

	return &Config{
		Port:                os.Getenv("PORT"),
		DbUser:              os.Getenv("DB_USER"),
		DbPassword:          os.Getenv("DB_PASSWORD"),
		DbName:              os.Getenv("DB_NAME"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Host:                os.Getenv("HOST"),
		JwtSecret:           os.Getenv("JWT_SECRET"),
		JwtExpiryMins:       jwtExpiryMins,
		Algorithm:           os.Getenv("ALGORITHM"),
		CloudflareAccountID: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareAPIToken:  os.Getenv("CLOUDFLARE_API_TOKEN"),
	}, nil
}
