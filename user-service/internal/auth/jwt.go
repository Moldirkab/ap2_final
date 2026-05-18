package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("super_secret_key")

func GenerateToken(userID string, role string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"role":    role,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
