package main

import (
	"carbon/models"
	"carbon/utilities"
	"carbon/validators"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("There was a problem loading the environment variables")
	}

	database := utilities.GetDatabaseObject()

	database.AutoMigrate(&models.User{})

	app := echo.New()
	app.Validator = &utilities.RequestValidator{Validator: validator.New()}

	app.POST("/auth/signup", validators.SignupValidator)
	app.POST("/auth/login", validators.LoginValidator)
	app.PUT("/users/:uuid/verify", validators.VerifyUserValidator)
	app.Logger.Fatal(app.Start(":8000"))
}
