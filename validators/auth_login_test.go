package validators

import (
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	godotenv.Load("../.env")
}

func TestLoginWithInvalidInputs(t *testing.T) {
	loginDto := new(dto.UserLoginDto)
	loginDtoJson, _ := json.Marshal(loginDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/login", loginDtoJson)
	err := LoginValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		response := parsedError.Message.(map[string][]chequerutilities.RequestError)

		assert.Equal(t, 2, len(response["errors"]))
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestLoginWithIncorrectInputs(t *testing.T) {
	loginDto := new(dto.UserLoginDto)

	if err := faker.FakeData(loginDto); err != nil {
		t.Fatal("There was a problem generating the fake data")
	}

	loginDtoJson, _ := json.Marshal(loginDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/login", loginDtoJson)

	err := LoginValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		assert.Equal(t, http.StatusBadRequest, parsedError.Code)

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, response["error"], "AUTH_002")
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestLoginSuccessfully(t *testing.T) {
	loginDto := new(dto.UserLoginDto)

	email := strings.ToLower(faker.Email())
	password := faker.Password()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := models.User{FirstName: faker.FirstName(), LastName: faker.LastName(), Email: email, Password: string(hashedPassword), AuthProvider: 1}

	database := chequerutilities.GetDatabaseObject()
	result := database.Save(&user)

	if result.Error != nil {
		t.Fatal("There was a problem creating the test user")
	}

	loginDto.Email = email
	loginDto.Password = password
	loginDtoJson, _ := json.Marshal(loginDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/login", loginDtoJson)

	err := LoginValidator(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}

		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NotNil(t, response["token"])

		userJson, _ := json.Marshal(response["user"])

		var user models.User

		json.Unmarshal(userJson, &user)

		assert.Equal(t, user.Email, loginDto.Email)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}
