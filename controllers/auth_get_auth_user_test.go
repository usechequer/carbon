package controllers

import (
	"carbon/models"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load("../.env")
}

func TestGetAuthUser(t *testing.T) {
	context, recorder := chequerutilities.GetTestUtilities(http.MethodGet, "/auth/me", []byte(""))

	var user models.User

	database := chequerutilities.GetDatabaseObject()

	result := database.Order("RAND()").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatal("There was a problem querying for the test user")
	}

	context.Set("user", user)
	err := GetAuthUser(context)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseUser models.User

		json.Unmarshal(recorder.Body.Bytes(), &responseUser)

		assert.Equal(t, user.Uuid.String(), responseUser.Uuid.String())
		assert.Equal(t, user.FirstName, responseUser.FirstName)
		assert.Equal(t, user.LastName, responseUser.LastName)
		assert.Equal(t, user.Email, responseUser.Email)
	} else {
		t.Fatal("The function completed wrongly with an error")
	}
}
