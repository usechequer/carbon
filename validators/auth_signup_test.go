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

func TestSignupValidatorInvalidInputs(t *testing.T) {
	signupDto := new(dto.UserSignupDto)
	signupDtoJson, _ := json.Marshal(signupDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/signup", signupDtoJson)

	err := SignupValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")
		assert.Equal(t, http.StatusBadRequest, parsedError.Code)

		var response = parsedError.Message.(map[string][]chequerutilities.RequestError)

		assert.Equal(t, 4, len(response["errors"]))
	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestSignupWithTakenEmail(t *testing.T) {
	signupDto := new(dto.UserSignupDto)

	if err := faker.FakeData(signupDto); err != nil {
		t.Fatal("There was a problem generating the fake data")
	}

	user := models.User{FirstName: signupDto.FirstName, LastName: signupDto.LastName, Email: strings.ToLower(signupDto.Email), Password: signupDto.Password, AuthProvider: 1}

	database := chequerutilities.GetDatabaseObject()

	result := database.Create(&user)

	if result.Error != nil {
		t.Fatal("There was an issue creating the test user")
	}

	newSignupDto := new(dto.UserSignupDto)

	if err := faker.FakeData(newSignupDto); err != nil {
		t.Fatal("There was a problem generating the fake data")
	}

	newSignupDto.Email = user.Email

	newSignupDtoJson, _ := json.Marshal(newSignupDto)

	context, _ := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/signup", newSignupDtoJson)

	err := SignupValidator(context)

	if assert.Error(t, err) {
		parsedError, ok := err.(*echo.HTTPError)

		assert.True(t, ok, "Expected an echo http error")

		assert.Equal(t, http.StatusBadRequest, parsedError.Code)

		response := parsedError.Message.(map[string]string)

		assert.Equal(t, "AUTH_001", response["error"])

	} else {
		t.Fatal("The function completed wrongly without an error")
	}
}

func TestSignUpValidatorSuccessful(t *testing.T) {
	signupDto := new(dto.UserSignupDto)
	err := faker.FakeData(signupDto)

	if err != nil {
		t.Fatal("There was a problem generating the fake data")
	}

	signupDtoJson, _ := json.Marshal(signupDto)

	context, recorder := chequerutilities.GetTestUtilities(http.MethodPost, "/auth/signup", signupDtoJson)

	if assert.NoError(t, SignupValidator(context)) {
		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NotNil(t, response["token"])

		userBytes, _ := json.Marshal(response["user"])

		var responseUser models.User
		json.Unmarshal(userBytes, &responseUser)

		database := chequerutilities.GetDatabaseObject()
		var user models.User

		database.Where("uuid = ?", responseUser.Uuid).First(&user)

		assert.Equal(t, user.FirstName, responseUser.FirstName)
		assert.Equal(t, user.LastName, responseUser.LastName)
		assert.Equal(t, user.Email, responseUser.Email)
	} else {
		t.Fatal("The function wrongly returned an error")
	}
}
