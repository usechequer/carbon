package utilities

import "github.com/labstack/echo/v4"

type Exception struct {
	StatusCode int
	Error      string
	Message    string
}

func ThrowException(context echo.Context, exception *Exception) error {
	return context.JSON(exception.StatusCode, map[string]string{"error": exception.Error, "message": exception.Message})
}
