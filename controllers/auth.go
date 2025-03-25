package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"math/rand"
	"net/http"
	"os"
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

	query := context.Request().URL.Query()
	query.Add("provider", provider)
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
		authProvider = 1
	case "github":
		authProvider = 2
	default:
		authProvider = 1
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
		user := models.User{FirstName: providerUser.FirstName, LastName: providerUser.LastName, Email: providerUser.Email, Password: GenerateRandomString(120), EmailVerifiedAt: GetTimestampPointer(time.Now()), AuthProvider: authProvider, Avatar: &providerUser.AvatarURL}
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

func GenerateRandomString(length int) string {
	const CHARSET = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	for i := range bytes {
		bytes[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}

	return string(bytes)
}
