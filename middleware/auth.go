package middleware

import (
	"carbon/models"
	"carbon/utilities"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		Authorization := context.Request().Header.Get("Authorization")

		authHeaderSplits := strings.Split(Authorization, " ")

		if len(authHeaderSplits) != 2 {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusUnauthorized, Message: "Not authenticated", Error: "AUTH_004"})
		}

		token := authHeaderSplits[1]

		decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !decodedToken.Valid {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusUnauthorized, Message: "Not authenticated", Error: "AUTH_004"})
		}

		userUuid, _ := decodedToken.Claims.GetSubject()

		var user models.User

		database := utilities.GetDatabaseObject()

		result := database.Where("uuid = ?", userUuid).First(&user)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusUnauthorized, Message: "Not authenticated", Error: "AUTH_004"})
		}

		context.Set("user", user)

		return next(context)
	}
}
