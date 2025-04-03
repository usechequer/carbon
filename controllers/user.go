package controllers

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	cloudinaryV2 "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

func UpdateUserAvatar(ctx echo.Context) error {
	cloudinary, err := cloudinaryV2.New()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "There was an issue processing the image upload"})
	}

	avatar, err := ctx.FormFile("avatar")

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "MALFORMED_REQUEST", "message": err.Error()})
	}

	src, err := avatar.Open()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "MALFORMED_REQUEST", "message": err.Error()})
	}

	defer src.Close()

	user := ctx.Get("user").(models.User)

	uploadResult, err := cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{Folder: fmt.Sprintf("%s/avatars", os.Getenv("CLOUDINARY_FOLDER")), ResourceType: "image", PublicID: user.Uuid.String(), Overwrite: api.Bool(true)})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "There was an issue processing the image upload"})
	}

	user.Avatar = &uploadResult.SecureURL

	database := utilities.GetDatabaseObject()

	database.Save(&user)

	return ctx.JSON(200, user)
}
