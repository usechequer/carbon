package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func getTimestampPointer(val time.Time) *time.Time {
	return &val
}

func VerifyUser(context echo.Context) error {
	user := context.Get("user").(models.User)

	database := utilities.GetDatabaseObject()

	user.EmailVerifiedAt = getTimestampPointer(time.Now())

	database.Save(&user)

	return context.JSON(http.StatusOK, user)
}

func UpdateUser(context echo.Context) error {
	user := context.Get("user").(models.User)
	updateUserDto := context.Get("updateUserDto").(*dto.UpdateUserDto)

	user.FirstName = updateUserDto.FirstName
	user.LastName = updateUserDto.LastName

	database := utilities.GetDatabaseObject()

	database.Save(&user)

	return context.JSON(http.StatusOK, user)
}
