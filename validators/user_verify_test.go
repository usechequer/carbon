package validators

import (
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
)

func TestVerifyUserWithInvalidUuid(t *testing.T) {
	verifyUserDto := new(dto.VerifyUserDto)
	verifyUserDtoJson, _ := json.Marshal(verifyUserDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPut, fmt.Sprintf("/users/%s/verify", verifyUserDto.Uuid.String()), verifyUserDtoJson)

	err := VerifyUserValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, "USER_001", response["error"])
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestVerifyUserWhoIsVerifiedAlready(t *testing.T) {
	getTimestampPointer := func(value time.Time) *time.Time {
		return &value
	}

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Where("email_verified_at IS NULL").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	user.EmailVerifiedAt = getTimestampPointer(time.Now())
	database.Save(&user)

	verifyUserDto := &dto.VerifyUserDto{Uuid: user.Uuid}
	verifyUserDtoJson, _ := json.Marshal(verifyUserDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPut, fmt.Sprintf("/users/%s/verify", user.Uuid.String()), verifyUserDtoJson)

	err := VerifyUserValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected echo http error")

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, "USER_002", response["error"])
	} else {
		t.Fatal("The function wrongly completed without an error")
	}
}

func TestVerifyUserSuccessfully(t *testing.T) {
	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Where("email_verified_at IS NULL").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	verifyUserDto := &dto.VerifyUserDto{Uuid: user.Uuid}
	verifyUserDtoJson, _ := json.Marshal(verifyUserDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPut, fmt.Sprintf("/users/%s/verify", user.Uuid.String()), verifyUserDtoJson)

	err := VerifyUserValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseUser models.User

		json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, user.Uuid.String(), responseUser.Uuid.String())
		assert.NotNil(t, responseUser.EmailVerifiedAt)

		var updatedUser models.User

		database.Where("id = ?", user.ID).First(&updatedUser)

		assert.NotNil(t, updatedUser.EmailVerifiedAt)
	} else {
		t.Fatal("The function wrongly returned an error.")
	}
}
