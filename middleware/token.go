package middleware

import (
	"carbon/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
)

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		token := context.Get("token").(chequerutilities.Token)

		userUuid := token.Subject

		var user models.User

		database := chequerutilities.GetDatabaseObject()

		result := database.Where("uuid = ?", userUuid).First(&user)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusUnauthorized, Message: "Not authenticated", Error: "AUTH_004"})
		}

		context.Set("user", user)

		return next(context)
	}
}
