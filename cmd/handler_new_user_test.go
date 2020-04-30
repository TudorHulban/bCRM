package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/steinfletcher/apitest"
	_ "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func Test1CreateUser(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointNewUser, NewUser)

	if assert.Nil(t, commons.CheckPgDB(e.Logger), "TEST - Could not connect to DB.") {
		f := make(url.Values)
		f.Set("teamid", "1")
		f.Set("code", commons.UXNano())
		f.Set("pass", "abcd")

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, commons.EndpointNewUser, strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		e.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}
}

func Test2CreateUser(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointNewUser, NewUser)

	if assert.Nil(t, commons.CheckPgDB(e.Logger), "TEST - Could not connect to DB.") {
		apitest.New().
			Handler(e).
			Method(http.MethodPost).
			URL(commons.EndpointNewUser).
			FormData("teamid", "1").
			FormData("code", commons.UXNano()).
			FormData("pass", "abcd").
			Expect(t).
			Status(http.StatusCreated).
			End()
	}
}
