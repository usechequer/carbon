package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"net/http"

	"github.com/labstack/echo/v4"
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

	return context.JSON(http.StatusCreated, user)
}

func Login(context echo.Context) error {
	user := context.Get("user").(models.User)

	token, err := utilities.GenerateJwtToken(user.Uuid.String())

	if err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusInternalServerError, Error: "AUTH_003", Message: "There was a problem generating the token."})
	}

	return context.JSON(http.StatusOK, map[string]interface{}{"token": token, "user": user})
}
