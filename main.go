package main

import (
	"carbon/controllers"
	"carbon/models"
	"carbon/utilities"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("There was a problem loading the environment variables")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("There was a problem connecting to the database")
	}

	database.AutoMigrate(&models.User{})

	app := echo.New()
	app.Validator = &utilities.RequestValidator{Validator: validator.New()}

	app.POST("/auth/signup", controllers.Signup)
	app.Logger.Fatal(app.Start(":8000"))
}
