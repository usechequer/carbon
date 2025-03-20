package validators

import (
	"carbon/controllers"
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func VerifyUserValidator(context echo.Context) error {
	verifyUserDto := new(dto.VerifyUserDto)

	if err := context.Bind(verifyUserDto); err != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "MALFORMED_REQUEST", Message: err.Error()})
	}

	var user models.User

	database := utilities.GetDatabaseObject()

	result := database.Where("uuid = ?", verifyUserDto.Uuid).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusNotFound, Error: "USER_001", Message: fmt.Sprintf("User with uuid %s does not exist", verifyUserDto.Uuid)})
	}

	if user.EmailVerifiedAt != nil {
		return utilities.ThrowException(context, &utilities.Exception{StatusCode: http.StatusBadRequest, Error: "USER_002", Message: fmt.Sprintf("User with uuid %s is verified already", verifyUserDto.Uuid)})
	}

	context.Set("user", user)

	return controllers.VerifyUser(context)
}
