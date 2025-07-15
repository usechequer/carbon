package validators

import (
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
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

	email := strings.ToLower(faker.Email())
	user := models.User{FirstName: faker.FirstName(), LastName: faker.LastName(), Email: email, Password: faker.Password(), AuthProvider: 1}

	database := chequerutilities.GetDatabaseObject()

	result := database.Save(&user)

	if result.Error != nil {
		t.Fatal("There was a problem creating the test user")
	}

	resetPasswordDto.Email = email

	resetPasswordDtoJson, _ := json.Marshal(resetPasswordDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password", resetPasswordDtoJson)

	err := ResetPasswordValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}
