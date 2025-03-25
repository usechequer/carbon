package main

import (
	"carbon/controllers"
	"carbon/middleware"
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

	utilities.RegisterOauthProviders()

	app := echo.New()
	app.Validator = &utilities.RequestValidator{Validator: validator.New()}

	app.POST("/auth/signup", validators.SignupValidator)
	app.POST("/auth/login", validators.LoginValidator)
	app.GET("/auth/:provider", controllers.OauthRedirectHandler)
	app.POST("/auth/:provider/callback", validators.OauthCallbackValidator)
	app.GET("/auth/:provider/callback", validators.OauthCallbackValidator)

	app.PUT("/users/:uuid/verify", validators.VerifyUserValidator)

	group := app.Group("/users/:uuid")
	group.Use(middleware.AuthMiddleware)
	group.PUT("", validators.UpdateUserValidator)

	app.Logger.Fatal(app.Start(":8000"))
}
