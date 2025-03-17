package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserSignupDto struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func Signup(context echo.Context) error {
	signupDto := new(UserSignupDto)

	if err := context.Bind(signupDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(signupDto); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, signupDto)
}
