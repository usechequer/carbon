package controllers

import (
	"carbon/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	chequerutilities "github.com/usechequer/utilities"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load("../.env")
}

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "test" {
		fmt.Println("=========== Running the migrations ============")

		database := chequerutilities.GetDatabaseObject()
		database.AutoMigrate(&models.User{})

		fmt.Println("=========== Migrations completed ============")

		fmt.Println("=========== Running the users seeder ===========")

		users := make([]models.User, 20)

		for i := 0; i < 20; i++ {
			users[i] = models.User{FirstName: faker.FirstName(), LastName: faker.LastName(), Email: strings.ToLower(faker.Email()), AuthProvider: 1, Password: faker.Password()}
		}

		database.Create(&users)

		fmt.Println("=========== Users seeder completed ===========")
	}

	m.Run()
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
