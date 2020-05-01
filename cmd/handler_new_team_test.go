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

func Test1_CreateTeam(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointNewTeam, NewTeam)

	if assert.Nil(t, commons.CheckPgDB(e.Logger), "TEST - Could not connect to DB.") {
		apitest.New().
			Handler(e).
			Method(http.MethodPost).
			URL(commons.EndpointNewTeam).
			FormData(commons.NewTeamFormName, "BLUE"+commons.UXSecs()).
			FormData(commons.NewTeamFormDescription, "DESC"+commons.UXSecs()).
			FormData(commons.NewTeamFormCODE, "CODE"+commons.UXSecs()).
			Expect(t).
			Status(http.StatusCreated).
			Assert(jsonpath.Matches(`$.ID`, `^\d+$`)).
			End()
	}
}
