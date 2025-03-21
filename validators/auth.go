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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignupValidator(context echo.Context) error {
	signupDto := new(dto.UserSignupDto)

	if err := context.Bind(signupDto); err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
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

func LoginValidator(context echo.Context) error {
	loginDto := new(dto.UserLoginDto)

	if err := context.Bind(loginDto); err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	if err := context.Validate(loginDto); err != nil {
		return err
	}

	var user models.User

	database := utilities.GetDatabaseObject()

	loginDto.Email = strings.ToLower(loginDto.Email)

	result := database.Where("email = ?", loginDto.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User does not exist with the specified email and password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))

	if err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User does not exist with the specified email and password"})
	}

	context.Set("user", user)

	return controllers.Login(context)
}
