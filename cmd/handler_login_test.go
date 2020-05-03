package main

import (
	"net/http"
	"testing"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func Test1_Login(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointLogin, LoginWithPassword)

	if assert.Nil(t, commons.CheckPgDB(e.Logger), "TEST - Could not connect to DB.") {
		apitest.New().
			Handler(e).
			Method(http.MethodPost).
			URL(commons.EndpointLogin).
			FormData(commons.LoginFormUserCode, "ADMIN").
			FormData(commons.LoginFormPass, "1234").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Matches(`$.ID`, `^\d+$`)).
			Assert(jsonpath.Present("$.token")). // https://play.golang.org/p/chl9Y9J8QVk
			End()
	}
}
