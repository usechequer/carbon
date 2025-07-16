package validators

import (
	"carbon/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
)

func TestUpdateUserWithInvalidUuid(t *testing.T) {
	uuid := faker.UUIDHyphenated()

	context, _ := chequerutilities.GetTestUtilities(http.MethodPut, fmt.Sprintf("/users/%s", uuid))

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	database.Order("RAND()").First(&user)

	context.Set("user", user)

	err := UpdateUserValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, http.StatusUnauthorized, parsedError.Code)
		assert.Equal(t, "AUTH_004", response["error"])
	} else {
		t.Fatal("The function wrongly completed without an error")
	}
}

func TestUpdateUserSuccessfutlly(t *testing.T) {

}
