package validators

import (
	"carbon/controllers"
	"carbon/models"
	"carbon/utilities"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func OauthCallbackValidator(context echo.Context) error {
	provider := context.Param("provider")
	query := context.Request().URL.Query()
	query.Add("provider", provider)
	context.Request().URL.RawQuery = query.Encode()

	request := context.Request()
	response := context.Response().Writer

	state := context.QueryParam("state")
	isLogin, _ := strconv.ParseBool(strings.Split(state, "=")[1])

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
