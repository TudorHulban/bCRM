package httphandlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// loginWithUserPassword is handler to perform user and password authentication against persisted data.
func loginWithUserPassword(c echo.Context) error {
	e := httpError{TheError: "credentials not found"}

	errValidateForm := validateForm(c)
	if errValidateForm != nil {
		e.TheError = "credentials could not be parsed"
		return c.JSON(http.StatusBadRequest, e)
	}
	// error is user info not found, we are not hiding this with 404
	u, errGetUser := blog.GetUserByCodeUnauthorized(c.FormValue("logincode"))
	log.Println("fetched:", u)
	if errGetUser != nil {
		return c.JSON(http.StatusForbidden, e)
	}
	// at this stage user info is found, to check password. error is password does not match not found, we are not hiding this with 404
	log.Println("Form user:", c.FormValue("logincode"), "Form pwd:", c.FormValue("password"), "PasswordSALT:", u.PasswordSALT)
	if !checkPasswordHash(c.FormValue("password"), u.PasswordSALT, u.PasswordHASH) {
		return c.JSON(http.StatusForbidden, e)
	}
	t, errJWT := createJWT()
	if errJWT != nil {
		return c.String(http.StatusInternalServerError, errJWT.Error())
	}
	log.Println(echo.Map{"token": t})
	return c.JSON(http.StatusOK, echo.Map{"token": t, "id": u.ID})
}
