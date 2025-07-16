package controllers

import (
	"carbon/models"
	"carbon/utilities"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	chequerutilities "github.com/usechequer/utilities"
)

func OauthRedirectHandler(context echo.Context) error {
	provider := context.Param("provider")

	isLoginStr := context.QueryParam("isLogin")
	var isLogin bool

	if isLoginStr == "" {
		isLogin = true
	} else {
		isLogin, _ = strconv.ParseBool(context.QueryParam("isLogin"))
	}

	query := context.Request().URL.Query()
	query.Add("provider", provider)
	query.Add("state", fmt.Sprintf("isLogin=%s", strconv.FormatBool(isLogin)))
	context.Request().URL.RawQuery = query.Encode()

	request := context.Request()
	response := context.Response().Writer

	gothic.Store = utilities.GetOauthSessionStore()

	if gothUser, err := gothic.CompleteUserAuth(response, request); err == nil {
		return context.JSON(http.StatusOK, gothUser)
	}

	gothic.BeginAuthHandler(response, request)
	return nil
}

func OauthCallbackHandler(context echo.Context) error {
	provider := context.Param("provider")
	var authProvider uint

	switch provider {
	case "google":
		authProvider = 2
	case "github":
		authProvider = 3
	default:
		authProvider = 2
	}

	isLogin := context.Get("isLogin").(bool)
	contextUser := context.Get("user")

	var user models.User
	var jwtToken string

	if isLogin {
		user = contextUser.(models.User)
		token, err := chequerutilities.GenerateJwtToken(user.Uuid.String())

		if err != nil {
			return chequerutilities.ThrowException(&chequerutilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
		}

		jwtToken = token
	} else {
		providerUser := contextUser.(goth.User)
		firstName, lastName := getOauthNames(providerUser.RawData)
		user := models.User{FirstName: firstName, LastName: lastName, Email: providerUser.Email, Password: utilities.GenerateRandomString(120), EmailVerifiedAt: getTimestampPointer(time.Now()), AuthProvider: authProvider, Avatar: &providerUser.AvatarURL}
		database := chequerutilities.GetDatabaseObject()
		result := database.Create(&user)

		if result.Error != nil {
			return context.JSON(http.StatusInternalServerError, map[string]string{"message": "There was a problem signing the user up."})
		}

		token, err := chequerutilities.GenerateJwtToken(user.Uuid.String())

		if err != nil {
			return chequerutilities.ThrowException(&chequerutilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
		}

		jwtToken = token
	}

	cookieToken := new(http.Cookie)
	cookieToken.Name = "token"
	cookieToken.Value = jwtToken
	cookieToken.Expires = time.Now().Add(time.Hour * 72)
	cookieToken.Path = "/"
	http.SetCookie(context.Response().Writer, cookieToken)

	return context.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
}

func getOauthNames(data map[string]interface{}) (firstName string, lastName string) {
	name := data["name"].(string)
	nameSplits := strings.Split(name, " ")
	firstName = nameSplits[0]
	lastName = strings.Join(nameSplits[1:], " ")

	return firstName, lastName
}
