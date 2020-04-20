package main

import (
	_ "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test1OK(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.GET(constants.EndpointNewUser, httphandlers.Live)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, constants.EndpointNewUser, nil)

	e.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
