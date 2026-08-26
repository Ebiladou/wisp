package utils

import (
	"errors"
	"time"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenType  = "ACCESS"
	RefreshTokenType = "REFRESH"

	RefreshTokenExpiry = 7 * 24 * time.Hour
)

type JWTClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	TokenID   string `json:"token_id,omitempty"`

	jwt.RegisteredClaims
}

func GenerateJWT(userID string, tokenType string, tokenID string, applicationConfig *config.Config, expiration time.Duration) (string, error) {

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

func ValidateJWT(tokenString string, applicationConfig *config.Config) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
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
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
