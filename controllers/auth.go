package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func getPasswordResetPointer(token string) *datatypes.JSON {
	passwordReset := datatypes.JSON([]byte(fmt.Sprintf(`{"token": "%s", "expires_at": "%s"}`, token, time.Now().Add(15*time.Minute).Format(time.RFC3339))))
	return &passwordReset
}

func GetAuthUser(context echo.Context) error {
	user := context.Get("user")

	return context.JSON(http.StatusOK, user)
}

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

func ResetPassword(context echo.Context) error {
	user := context.Get("user").(models.User)

	token := generateRandomString(100)
	user.PasswordReset = getPasswordResetPointer(token)

	database := utilities.GetDatabaseObject()

	database.Save(&user)

	return context.JSON(http.StatusOK, map[string]string{"message": "Reset password email sent successfully"})
}

func ConfirmResetPassword(context echo.Context) error {
	user := context.Get("user").(models.User)
	confirmResetPasswordDto := context.Get("confirmResetPasswordDto").(*dto.ConfirmResetPasswordDto)

	password, _ := bcrypt.GenerateFromPassword([]byte(confirmResetPasswordDto.Password), 14)

	user.Password = string(password)
	user.PasswordReset = nil

	database := utilities.GetDatabaseObject()
	database.Save(&user)

	return context.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
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
