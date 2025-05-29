package validators

import (
	"carbon/controllers"
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	chequerutilities "github.com/usechequer/utilities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SignupValidator(context echo.Context) error {
	signupDto := new(dto.UserSignupDto)

	if err := context.Bind(signupDto); err != nil {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	if err := context.Validate(signupDto); err != nil {
		// return err
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: "akak"})
	}

	database := chequerutilities.GetDatabaseObject()

	var user models.User

	signupDto.Email = strings.ToLower(signupDto.Email)

	result := database.Where("email = ?", signupDto.Email).First(&user)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_001", Message: fmt.Sprintf("User with email address %s exists already", signupDto.Email)})
	}

	context.Set("signupDto", signupDto)

	return controllers.Signup(context)
}

func LoginValidator(context echo.Context) error {
	loginDto := new(dto.UserLoginDto)

	if err := context.Bind(loginDto); err != nil {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	if err := context.Validate(loginDto); err != nil {
		return err
	}

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	loginDto.Email = strings.ToLower(loginDto.Email)

	result := database.Where("email = ?", loginDto.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User does not exist with the specified email and password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))

	if err != nil {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User does not exist with the specified email and password"})
	}

	context.Set("user", user)

	return controllers.Login(context)
}

func ResetPasswordValidator(context echo.Context) error {
	resetPasswordDto := new(dto.ResetPasswordDto)

	if err := context.Bind(resetPasswordDto); err != nil {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	if err := context.Validate(resetPasswordDto); err != nil {
		return err
	}

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Where("email = ?", resetPasswordDto.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusNotFound, Error: "USER_001", Message: fmt.Sprintf("User with email %s does not exist", resetPasswordDto.Email)})
	}

	context.Set("user", user)

	return controllers.ResetPassword(context)
}

func ConfirmResetPasswordValidator(context echo.Context) error {
	confirmResetPasswordDto := new(dto.ConfirmResetPasswordDto)

	if err := context.Bind(confirmResetPasswordDto); err != nil {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	if err := context.Validate(confirmResetPasswordDto); err != nil {
		return err
	}

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.First(&user, datatypes.JSONQuery("password_reset").Equals(confirmResetPasswordDto.Token, "token"))

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusNotFound, Error: "USER_001", Message: fmt.Sprintf("User with password reset token %s does not exist", confirmResetPasswordDto.Token)})
	}

	var pwReset map[string]interface{}

	json.Unmarshal(*user.PasswordReset, &pwReset)

	expiresAt, _ := time.Parse(time.RFC3339, pwReset["expires_at"].(string))

	if expiresAt.Compare(time.Now()) < 0 {
		return chequerutilities.ThrowException(context, &chequerutilities.Exception{StatusCode: http.StatusBadRequest, Error: "USER_003", Message: "Password reset token has expired"})
	}

	context.Set("user", user)
	context.Set("confirmResetPasswordDto", confirmResetPasswordDto)

	return controllers.ConfirmResetPassword(context)
}
