package main

import (
	"net/http"

	"github.com/TudorHulban/bCRM/models"
	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// LoginWithPassword is handler to perform user and password authentication against persisted data.
func LoginWithPassword(c echo.Context) error {
	var e httpError
	c.Logger().Debug("Login w Password")

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

	// error is user info not found, we are not hiding this with 404
	m, errGetUser := models.NewUser(c, u, true)
	user, errGetUser := m.GetUserByCodeUnauthorized()
	if errGetUser != nil {
		e.TheError = errGetUser.Error()
		return c.JSON(http.StatusNotFound, e)
	}
	log.Debug("fetched user:", user)

	// at this stage user info is found, to check password. error is password does not match not found, we are not hiding this with 404
	if !checkPasswordHash(c.FormValue(u.LoginPWD), user.PasswordSALT, user.PasswordHASH) {
		e.TheError = "bad credentials"
		return c.JSON(http.StatusForbidden, e)
	}

	// user is authenticated, no authorization for now
	t, errJWT := createJWT(commons.TokenExpirationSeconds)
	if errJWT != nil {
		e.TheError = errJWT.Error()
		return c.JSON(http.StatusInternalServerError, e)
	}
	log.Debug(echo.Map{"token": t})
	return c.JSON(http.StatusOK, echo.Map{"token": t, "id": user.ID})
}
