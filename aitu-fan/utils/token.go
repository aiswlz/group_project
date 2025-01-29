package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/group-project/aitu-fan/aitu-fan/models"
)

var secretKey = []byte("your_secret_key")

func ValidateToken(tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		role := claims["role"].(string)

		return &models.User{
			Username: username,
			Role:     role,
		}, nil
	}

	return nil, errors.New("Invalid token claims")
}
