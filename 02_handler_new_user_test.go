package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test2NewUser(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.GET(EndpointNewUser, NewUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, EndpointNewUser, nil)

	e.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
