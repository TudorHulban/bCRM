package httphandlers_test

import (
	_ "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/TudorHulban/bCRM/pkg/httphandlers"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test1OK(t *testing.T) {
	e := echo.New()
	e.GET(constants.EndpointNewUser, httphandlers.HandlerOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, constants.EndpointNewUser, nil)

	e.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
