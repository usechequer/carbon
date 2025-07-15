package validators

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/datatypes"
)

func TestConfirmResetPasswordWithInvalidInputs(t *testing.T) {
	confirmResetPasswordDto := new(dto.ConfirmResetPasswordDto)
	confirmResetPasswordDtoJson, _ := json.Marshal(confirmResetPasswordDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password/confirm", confirmResetPasswordDtoJson)

	err := ConfirmResetPasswordValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		response := parsedError.Message.(map[string][]chequerutilities.RequestError)

		assert.Equal(t, 2, len(response["errors"]))
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestConfirmResetPasswordWithIncorrectToken(t *testing.T) {
	confirmResetPasswordDto := new(dto.ConfirmResetPasswordDto)
	confirmResetPasswordDto.Token = faker.FirstName()
	confirmResetPasswordDto.Password = faker.Password()
	confirmResetPasswordDtoJson, _ := json.Marshal(confirmResetPasswordDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password/confirm", confirmResetPasswordDtoJson)

	err := ConfirmResetPasswordValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected echo http error")

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, "USER_001", response["error"])
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestConfirmResetPasswordWithExpiredToken(t *testing.T) {
	passwordResetToken := utilities.GenerateRandomString(50)
	passwordReset := datatypes.JSON([]byte(fmt.Sprintf(`{"token": "%s", "expires_at": "%s"}`, passwordResetToken, time.Now().Add(-1*time.Minute).Format(time.RFC3339))))

	user := models.User{FirstName: faker.FirstName(), LastName: faker.LastName(), Email: faker.Email(), Password: faker.Password(), AuthProvider: 1, PasswordReset: &passwordReset}

	database := chequerutilities.GetDatabaseObject()
	result := database.Save(&user)

	if result.Error != nil {
		t.Fatal("There was a problem creating the test user")
	}

	confirmResetPasswordDto := new(dto.ConfirmResetPasswordDto)
	confirmResetPasswordDto.Token = passwordResetToken
	confirmResetPasswordDto.Password = faker.Password()
	confirmResetPasswordDtoJson, _ := json.Marshal(confirmResetPasswordDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password/confirm", confirmResetPasswordDtoJson)

	err := ConfirmResetPasswordValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, "USER_003", response["error"])
	} else {
		t.Fatal("The function wrongly returned without an error")
	}

}
