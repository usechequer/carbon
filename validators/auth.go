package validators

import (
	"carbon/controllers"
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignupValidator(context echo.Context) error {
	signupDto := new(dto.UserSignupDto)

	if err := context.Bind(signupDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(signupDto); err != nil {
		return err
	}

	database := utilities.GetDatabaseObject()

	var user models.User

	signupDto.Email = strings.ToLower(signupDto.Email)

	result := database.Where("email = ?", signupDto.Email).First(&user)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_001", Message: fmt.Sprintf("User with email address %s exists already", signupDto.Email)})
	}

	context.Set("signupDto", signupDto)

	return controllers.Signup(context)
}
