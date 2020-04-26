package models

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

func isValidStruct(s interface{}, log echo.Logger) error {
	v := validator.New()

	errValid := v.Struct(s)
	if errValid != nil {
		log.Debugf("structure %s is invalid. error: %s", s, errValid)
		return errValid
	}
	return nil
}
