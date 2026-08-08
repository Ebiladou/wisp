package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ebiladou/wisp/internal/config"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string,
) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

func GenerateToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func HashToken(token string) string {
	return fmt.Sprintf("%x", bcryptHash([]byte(token)))
}

func bcryptHash(value []byte) []byte {
	hash, err := bcrypt.GenerateFromPassword(
		value,
		bcrypt.DefaultCost,
	)

	if err != nil {
		return []byte{}
	}

	return hash
}

func GenerateRandomTokenID() (string, error) {

	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

const (
	AccessTokenType  = "ACCESS"
	RefreshTokenType = "REFRESH"

	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

type JWTClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	TokenID   string `json:"token_id,omitempty"`

	jwt.RegisteredClaims
}

func GenerateJWT(
	userID string,
	tokenType string,
	tokenID string,
	applicationConfig *config.Config,
	expiration time.Duration,
) (string, error) {

	currentTime := time.Now()

	claims := JWTClaims{
		UserID:    userID,
		TokenType: tokenType,
		TokenID:   tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.New().String(),
			IssuedAt: jwt.NewNumericDate(
				currentTime,
			),
			ExpiresAt: jwt.NewNumericDate(
				currentTime.Add(expiration),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(applicationConfig.JwtSecret),
	)
}

func ValidateJWT(
	tokenString string,
	applicationConfig *config.Config,
) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"Unexpected signing method",
				)
			}

			return []byte(
				applicationConfig.JwtSecret,
			), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
