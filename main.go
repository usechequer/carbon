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

	app.POST("/auth/reset-password", validators.ResetPasswordValidator)
	app.POST("/auth/reset-password/confirm", validators.ConfirmResetPasswordValidator)

	app.GET("/auth/:provider", controllers.OauthRedirectHandler)
	app.GET("/auth/:provider/callback", validators.OauthCallbackValidator)

	app.PUT("/users/:uuid/verify", validators.VerifyUserValidator)

	authGroup := app.Group("/auth/me")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.GET("", controllers.GetAuthUser)

	userGroup := app.Group("/users/:uuid")
	userGroup.Use(middleware.AuthMiddleware)
	userGroup.PUT("", validators.UpdateUserValidator)

	app.Logger.Fatal(app.Start(":8000"))
}
