package main

import (
	"net/http"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
)

// LoginWithPassword is handler to perform user and password authentication against persisted data.
// RAW testing: curl -d "code=ADMIN&pass=1234" -X POST http://localhost:8001/login
func LoginWithPassword(c echo.Context) error {
	c.Logger().Debug("Login w Password")

	var e httpError
	if len(c.FormValue(commons.LoginFormUserCode)) == 0 {
		c.Logger().Debug("received ", commons.LoginFormUserCode, " as: ", c.FormValue(commons.LoginFormUserCode))
		e.TheError = commons.LoginFormUserCode + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	// password cannot be nil.
	if len(c.FormValue(commons.LoginFormPass)) == 0 {
		c.Logger().Debug("received ", commons.LoginFormPass, " as: ", c.FormValue(commons.LoginFormPass))
		e.TheError = commons.LoginFormPass + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var u models.UserFormData
	u.LoginCODE = c.FormValue(commons.LoginFormUserCode)
	u.LoginPWD = c.FormValue(commons.LoginFormPass)
	c.Logger().Debug("form data:", u)

	// error is user info not found, we are not hiding this with 404
	m, errGetUser := models.NewUser(c, u, true)
	user, errGetUser := m.GetUserByCodeUnauthorized()
	if errGetUser != nil {
		e.TheError = errGetUser.Error()
		return c.JSON(http.StatusNotFound, e)
	}

	// at this stage user info is found, to check password. error is password does not match not found, we are not hiding this with 404
	if !checkPasswordHash(u.LoginPWD, user.PasswordSALT, user.PasswordHASH) {
		e.TheError = "bad credentials"
		c.Logger().Debug(e.TheError, ": ", u.LoginPWD)
		return c.JSON(http.StatusForbidden, e)
	}

	// user is authenticated, no authorization for now
	sessionID := randomString(commons.SessionIDLength)
	// store session ID in cache for future requests
	getSessionIDCache().SetTTL(user.LoginCODE, sessionID, commons.SessionIDExpirationSeconds)

	return c.JSON(http.StatusOK, echo.Map{"sessionID": sessionID, "ID": user.ID})
}
