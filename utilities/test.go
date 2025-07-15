package utilities

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	chequerutilities "github.com/usechequer/utilities"
)

func GetTestUtilities(method string, path string, requestBodies ...[]byte) (context echo.Context, recorder *httptest.ResponseRecorder) {
	app := echo.New()
	app.Validator = &chequerutilities.RequestValidator{Validator: validator.New()}

	var request *http.Request

	if len(requestBodies) == 1 {
		request = httptest.NewRequest(method, path, strings.NewReader(string(requestBodies[0])))
	} else {
		request = httptest.NewRequest(method, path, strings.NewReader(""))
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	recorder = httptest.NewRecorder()
	context = app.NewContext(request, recorder)

	return
}
