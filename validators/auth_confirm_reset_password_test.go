package validators

import (
	"carbon/dto"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
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
