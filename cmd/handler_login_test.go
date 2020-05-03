package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/TudorHulban/bCRM/pkg/cache"

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

	// create cache
	theCache = cache.New(e.Logger)
	defer theCache.Close()

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
			Assert(jsonpath.Present("$.sessionID")). // https://play.golang.org/p/chl9Y9J8QVk
			End()
	}
}

func Test2_LoginTTL(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.POST(commons.EndpointLogin, LoginWithPassword)

	// create cache
	theCache = cache.New(e.Logger)
	defer theCache.Close()

	testLogin := func() string {
		f := make(url.Values)
		f.Set(commons.LoginFormUserCode, "ADMIN")
		f.Set(commons.LoginFormPass, "1234")

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, commons.EndpointLogin, strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		e.ServeHTTP(w, req)
		resp := w.Result()
		e.Logger.Debug("body: ", w.Body.String())
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data map[string]interface{}
		err := json.Unmarshal([]byte(w.Body.String()), &data)
		if !assert.Nil(t, err, "Unmarshal of response did not succeed.") {
			e.Logger.Fatal("could not unmarshal response")
		}

		e.Logger.Debug("result: ", data["sessionID"])
		return (data["sessionID"]).(string)
	}

	sessionID1 := testLogin()
	sessionID2 := testLogin()
	assert.EqualValues(t, sessionID1, sessionID2, "caching of session ID did not work", sessionID1, sessionID2)

	time.Sleep(time.Duration(commons.SessionIDExpirationSeconds+1) * time.Second)
	sessionID3 := testLogin()
	assert.NotEqual(t, sessionID1, sessionID3, "useed cached value instead of new due to TTL expired", sessionID1, sessionID3)
}
