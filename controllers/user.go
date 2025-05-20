package controllers

import (
	"carbon/dto"
	"carbon/models"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	cloudinaryV2 "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v4"
	chequerutilities "github.com/usechequer/utilities"
)

func getTimestampPointer(val time.Time) *time.Time {
	return &val
}

var avatarErrorMessage = "There was an issue processing the avatar upload"

func VerifyUser(context echo.Context) error {
	user := context.Get("user").(models.User)

	database := chequerutilities.GetDatabaseObject()

	user.EmailVerifiedAt = getTimestampPointer(time.Now())

	database.Save(&user)

	return context.JSON(http.StatusOK, user)
}

func UpdateUser(ctx echo.Context) error {
	user := ctx.Get("user").(models.User)
	updateUserDto := ctx.Get("updateUserDto").(*dto.UpdateUserDto)

	if len(updateUserDto.FirstName) > 0 {
		user.FirstName = updateUserDto.FirstName
	}

	if len(updateUserDto.LastName) > 0 {
		user.LastName = updateUserDto.LastName
	}

	if len(updateUserDto.CurrentProjectUuid) > 0 {
		user.CurrentProjectUuid = &updateUserDto.CurrentProjectUuid
	}

	database := chequerutilities.GetDatabaseObject()

	avatar, err := ctx.FormFile("avatar")

	if err != nil {
		database.Save(&user)
		return ctx.JSON(http.StatusOK, user)
	}

	avatarSrc, err := uploadAvatar(avatar, user.Uuid.String())

	if err != nil {
		return chequerutilities.ThrowException(ctx, &chequerutilities.Exception{Error: "USER_004", Message: avatarErrorMessage})
	}

	user.Avatar = &avatarSrc

	database.Save(&user)

	return ctx.JSON(http.StatusOK, user)
}

func uploadAvatar(avatar *multipart.FileHeader, userUuid string) (avatarSrc string, err error) {
	cloudinary, err := cloudinaryV2.New()

	if err != nil {
		return "", errors.New(avatarErrorMessage)
	}

	src, err := avatar.Open()

	if err != nil {
		return "", errors.New(avatarErrorMessage)
	}

	defer src.Close()

	uploadResult, err := cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{Folder: fmt.Sprintf("%s/avatars", os.Getenv("CLOUDINARY_FOLDER")), ResourceType: "image", PublicID: userUuid, Overwrite: api.Bool(true)})

	if err != nil {
		return "", errors.New(avatarErrorMessage)
	}

	return uploadResult.SecureURL, nil
}
