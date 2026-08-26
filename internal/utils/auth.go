package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
