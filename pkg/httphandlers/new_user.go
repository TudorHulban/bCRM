package httphandlers

import (
	"net/http"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/TudorHulban/bCRM/structs"
	"github.com/TudorHulban/bCRM/variables"
	"github.com/labstack/echo"
)

// NewUser Creates a new user. To be used by user management roles.
func NewUser(c echo.Context) error {
	var e httpError

	if len(c.FormValue(constants.NewUserFormName)) == 0 {
		e.TheError = constants.NewUserFormName + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(constants.NewUserFormUserCode)) == 0 {
		e.TheError = constants.NewUserFormUserCode + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	if len(c.FormValue(constants.NewUserFormPass)) == 0 {
		e.TheError = constants.NewUserFormPass + " information is not valid"
		return c.JSON(http.StatusNotAcceptable, e)
	}

	var co structs.Contact
	co.FirstName = c.FormValue(constants.NewUserFormName)
	co.CompanyEmail = c.FormValue(constants.NewUserFormEmail)

	var u structs.User
	u.LoginCODE = c.FormValue(constants.NewUserFormUserCode)
	u.LoginPWD = c.FormValue(constants.NewUserFormPass)
	u.ContactInfo = append(u.ContactInfo, &co)
	u.SecurityGroup = constants.SecuGrpUser

	errAdd := variables.GStore.New(u)
	if errAdd != nil {
		e.TheError = errAdd.Error()
		return c.JSON(http.StatusInternalServerError, e)
	}
	return c.JSON(http.StatusOK, u) // TODO: maybe we should send less info for user.
}
