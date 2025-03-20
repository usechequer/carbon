package utilities

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(subject string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"iss": os.Getenv("APP_NAME"),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return token, nil
}
