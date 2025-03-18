package validators

import (
	"carbon/controllers"
	"carbon/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SignupValidator(context echo.Context) error {
	signupDto := new(dto.UserSignupDto)

	if err := context.Bind(signupDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(signupDto); err != nil {
		return err
	}

	context.Set("signupDto", signupDto)

	return controllers.Signup(context)
}
