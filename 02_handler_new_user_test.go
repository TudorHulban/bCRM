package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test2NewUser(t *testing.T) {
	form := make(url.Values)
	form.Set("name", "John Smith")
	form.Set("code", "xxxx")
	form.Set("pass", "1234")

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	req := httptest.NewRequest(http.MethodPost, EndpointNewUser, strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	w := httptest.NewRecorder()
	echoCtx := e.NewContext(req, w)

	//log.Print(NewUser(echoCtx))
	assert.NoError(t, NewUser(echoCtx))
}
