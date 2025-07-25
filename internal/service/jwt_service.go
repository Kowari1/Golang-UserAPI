package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userID string, isAdmin bool, login string, key []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    isAdmin,
		"login":   login,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"jti":     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
