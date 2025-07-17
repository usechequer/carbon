package validators

import (
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
)

func TestResetPasswordWithInvalidInputs(t *testing.T) {
	resetPasswordDto := new(dto.ResetPasswordDto)
	resetPasswordDto.Email = faker.FirstName()
	resetPasswordDtoJson, _ := json.Marshal(resetPasswordDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password", resetPasswordDtoJson)

	err := ResetPasswordValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected echo http error")

		response := parsedError.Message.(map[string][]chequerutilities.RequestError)

		assert.Equal(t, 1, len(response["errors"]))
	} else {
		t.Fatal("The function wrongly completed without an error")
	}
}

func TestResetPasswordWithIncorrectEmail(t *testing.T) {
	resetPasswordDto := new(dto.ResetPasswordDto)
	resetPasswordDto.Email = faker.Email()
	resetPasswordDtoJson, _ := json.Marshal(resetPasswordDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password", resetPasswordDtoJson)

	err := ResetPasswordValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}

func TestResetPasswordSuccessfully(t *testing.T) {
	resetPasswordDto := new(dto.ResetPasswordDto)

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Order("RAND()").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	resetPasswordDto.Email = user.Email

	resetPasswordDtoJson, _ := json.Marshal(resetPasswordDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password", resetPasswordDtoJson)

	err := ResetPasswordValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}
