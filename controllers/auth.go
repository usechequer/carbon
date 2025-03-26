package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func Signup(context echo.Context) error {
	var signupDto *dto.UserSignupDto = context.Get("signupDto").(*dto.UserSignupDto)

	password, _ := bcrypt.GenerateFromPassword([]byte(signupDto.Password), 14)

	user := models.User{FirstName: signupDto.FirstName, LastName: signupDto.LastName, Email: signupDto.Email, Password: string(password), AuthProvider: 1}

	database := utilities.GetDatabaseObject()

	result := database.Create(&user)

	if result.Error != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "There was a problem signing the user up."})
	}

	token, err := utilities.GenerateJwtToken(user.Uuid.String())

	if err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
	}

	return context.JSON(http.StatusCreated, map[string]interface{}{"token": token, "user": user})
}

func Login(context echo.Context) error {
	user := context.Get("user").(models.User)

	token, err := utilities.GenerateJwtToken(user.Uuid.String())

	if err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
	}

	return context.JSON(http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

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
		token, err := utilities.GenerateJwtToken(user.Uuid.String())

		if err != nil {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
		}

		jwtToken = token
	} else {
		providerUser := contextUser.(goth.User)
		firstName, lastName := getOauthNames(providerUser.RawData)
		user := models.User{FirstName: firstName, LastName: lastName, Email: providerUser.Email, Password: generateRandomString(120), EmailVerifiedAt: GetTimestampPointer(time.Now()), AuthProvider: authProvider, Avatar: &providerUser.AvatarURL}
		database := utilities.GetDatabaseObject()
		result := database.Create(&user)

		if result.Error != nil {
			return context.JSON(http.StatusInternalServerError, map[string]string{"message": "There was a problem signing the user up."})
		}

		token, err := utilities.GenerateJwtToken(user.Uuid.String())

		if err != nil {
			return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
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

func generateRandomString(length int) string {
	const CHARSET = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	for i := range bytes {
		bytes[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}

	return string(bytes)
}
