package validators

import (
	"carbon/dto"
	"carbon/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
)

func TestSignUpValidator(t *testing.T) {
	godotenv.Load("../.env")
	signupDto := new(dto.UserSignupDto)
	err := faker.FakeData(&signupDto)

	if err != nil {
		t.Fatalf("There was an issue generating the fake data")
	}

	signupDtoJson, _ := json.Marshal(signupDto)

	app := echo.New()
	app.Validator = &chequerutilities.RequestValidator{Validator: validator.New()}
	request := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(string(signupDtoJson)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	recorder := httptest.NewRecorder()
	context := app.NewContext(request, recorder)

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
	}
}
