package httphandlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	a := assert.New(t)

	req, err := http.NewRequest("POST", constants.EndpointNewUser, nil)
	a.Nil(err)

	e := echo.New()

	respRecorder := httptest.NewRecorder()
	ectx := e.NewContext(req, respRecorder)

	a.NotNil(t, e.Router())

	e.DefaultHTTPErrorHandler(errors.New("error"), ectx)
	assert.Equal(t, http.StatusOK, respRecorder.Code)
}
