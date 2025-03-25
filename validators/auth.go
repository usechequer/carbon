package validators

import (
	"carbon/controllers"
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
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

func OauthCallbackValidator(context echo.Context) error {
	provider := context.Param("provider")
	query := context.Request().URL.Query()
	query.Add("provider", provider)
	context.Request().URL.RawQuery = query.Encode()

	request := context.Request()
	response := context.Response().Writer

	isLoginStr := context.QueryParam("isLogin")
	var isLogin bool

	if isLoginStr == "" {
		isLogin = true
	} else {
		isLogin, _ = strconv.ParseBool(context.QueryParam("isLogin"))
	}

	gothic.Store = utilities.GetOauthSessionStore()

	providerUser, err := gothic.CompleteUserAuth(response, request)

	if err != nil {
		fmt.Println(err)
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_005", Message: fmt.Sprintf("There was an issue retrieving user information from %s", provider)})
	}

	context.Set("isLogin", isLogin)

	var user models.User

	database := utilities.GetDatabaseObject()

	result := database.Where("email = ?", providerUser.Email).First(&user)

	if isLogin {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User does not exist with the specified email"})
		}

		context.Set("user", user)
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "AUTH_002", Message: "User exists already"})
		}

		context.Set("user", providerUser)
	}

	return controllers.OauthCallbackHandler(context)
}
