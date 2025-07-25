package validators

import (
	"carbon/dto"
	"carbon/models"
	"carbon/utilities"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

	var user models.User

	database := chequerutilities.GetDatabaseObject()
	result := database.Where("password_reset IS NULL").Order("RAND()").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	user.PasswordReset = &passwordReset
	database.Save(&user)

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

func TestConfirmResetPasswordSuccessfully(t *testing.T) {
	passwordResetToken := utilities.GenerateRandomString(50)
	passwordReset := datatypes.JSON([]byte(fmt.Sprintf(`{"token": "%s", "expires_at": "%s"}`, passwordResetToken, time.Now().Add(5*time.Minute).Format(time.RFC3339))))

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Where("password_reset IS NULL").Order("RAND()").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	user.PasswordReset = &passwordReset
	database.Save(&user)

	confirmResetPasswordDto := new(dto.ConfirmResetPasswordDto)
	confirmResetPasswordDto.Token = passwordResetToken
	confirmResetPasswordDto.Password = faker.Password()
	confirmResetPasswordDtoJson, _ := json.Marshal(confirmResetPasswordDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/reset-password/confirm", confirmResetPasswordDtoJson)

	err := ConfirmResetPasswordValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var updatedUser models.User
		database.Where("id = ?", user.ID).First(&updatedUser)

		assert.Nil(t, updatedUser.PasswordReset)

		err := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(confirmResetPasswordDto.Password))

		assert.Nil(t, err)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}
