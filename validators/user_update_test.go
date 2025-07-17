package validators

import (
	"bytes"
	"carbon/models"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
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

func TestUpdateUserSuccessfully(t *testing.T) {
	app := echo.New()
	app.Validator = &chequerutilities.RequestValidator{Validator: validator.New()}

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Order("RAND()").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	var request *http.Request

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	firstName := faker.FirstName()
	lastName := faker.LastName()
	currentProjectUuid := faker.UUIDHyphenated()

	writer.WriteField("first_name", firstName)
	writer.WriteField("last_name", lastName)
	writer.WriteField("current_project_uuid", currentProjectUuid)
	writer.WriteField("uuid", user.Uuid.String())

	multipart, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	multipart.Write([]byte("sample avatar"))
	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", user.Uuid.String()), body)
	request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	context := app.NewContext(request, recorder)
	context.Set("user", user)
	context.SetParamNames("uuid")
	context.SetParamValues(user.Uuid.String())

	err := UpdateUserValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseUser models.User

		json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, firstName, responseUser.FirstName)
		assert.Equal(t, lastName, responseUser.LastName)
		assert.Equal(t, currentProjectUuid, responseUser.CurrentProjectUuid.String())
		assert.Equal(t, "", *responseUser.Avatar)

		var updatedUser models.User

		database.Where("id = ?", user.ID).First(&updatedUser)

		assert.Equal(t, firstName, updatedUser.FirstName)
		assert.Equal(t, lastName, updatedUser.LastName)
		assert.Equal(t, currentProjectUuid, updatedUser.CurrentProjectUuid.String())
		assert.Equal(t, "", *updatedUser.Avatar)
	} else {
		t.Fatal("The function returned wrongly with an error")
	}
}
