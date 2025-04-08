package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() string {
	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "admin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenStr
}
