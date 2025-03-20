package controllers

import (
	"carbon/models"
	"carbon/utilities"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetTimestampPointer(val time.Time) *time.Time {
	return &val
}

func VerifyUser(context echo.Context) error {
	user := context.Get("user").(models.User)

	database := utilities.GetDatabaseObject()

	user.EmailVerifiedAt = GetTimestampPointer(time.Now())

	database.Save(&user)

	return context.JSON(http.StatusOK, user)
}
