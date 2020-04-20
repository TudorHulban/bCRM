package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test1Live(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.GET(EndpointLive, Live)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, EndpointLive, nil)

	e.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
