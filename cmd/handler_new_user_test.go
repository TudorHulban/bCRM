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
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func Test1CreateUser(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointNewUser, NewUser)

	if assert.Nil(t, commons.CheckPgDB(e.Logger), "TEST - Could not connect to DB.") {
		f := make(url.Values)
		f.Set(commons.NewUserFormTeamID, "1")
		f.Set(commons.NewUserFormUserCode, commons.UXNano())
		f.Set(commons.NewUserFormPass, "abcd")

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
			FormData(commons.NewUserFormTeamID, "1").
			FormData(commons.NewUserFormUserCode, commons.UXNano()).
			FormData(commons.NewUserFormPass, "abcd").
			Expect(t).
			Status(http.StatusCreated).
			Assert(jsonpath.Matches(`$.ID`, `^\d+$`)).
			End()
	}
}
