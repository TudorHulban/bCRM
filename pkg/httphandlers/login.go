package httphandlers

import (
	"net/http"

	"github.com/TudorHulban/bCRM/pkg/constants"
	"github.com/TudorHulban/bCRM/pkg/variables"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// LoginWithPassword is handler to perform user and password authentication against persisted data.
func LoginWithPassword(c echo.Context) error {
	e := httpError{TheError: "credentials not found"}

	userCode := c.FormValue(constants.LoginFormUser)
	if len(userCode) == 0 {
		e.TheError = "user code could not be parsed"
		return c.JSON(http.StatusBadRequest, e)
	}
	// error is user info not found, we are not hiding this with 404
	u, errGetUser := variables.GStore.GetUserByUserCode(userCode)
	if errGetUser != nil {
		e.TheError = errGetUser.Error()
		return c.JSON(http.StatusForbidden, e)
	}
	log.Debug("fetched user:", u)

	// at this stage user info is found, to check password. error is password does not match not found, we are not hiding this with 404
	if !checkPasswordHash(c.FormValue(constants.LoginFormPass), u.PasswordSALT, u.PasswordHASH) {
		e.TheError = "bad credentials"
		return c.JSON(http.StatusForbidden, e)
	}

	// user is authenticated, no authorization for now
	t, errJWT := createJWT(constants.TokenExpirationSeconds)
	if errJWT != nil {
		e.TheError = errJWT.Error()
		return c.JSON(http.StatusInternalServerError, e)
	}
	log.Debug(echo.Map{"token": t})
	return c.JSON(http.StatusOK, echo.Map{"token": t, "id": u.ID})
}
