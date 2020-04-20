package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

var formData = `{"name":"John Smith", "code":"xxxx", "pass":"1234"}`

func Test2NewUser(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	req := httptest.NewRequest(http.MethodPost, EndpointNewUser, strings.NewReader(formData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	w := httptest.NewRecorder()
	echoCtx := e.NewContext(req, w)

	assert.NoError(t, NewUser(echoCtx))
}
