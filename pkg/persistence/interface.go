package persistence

import (
	"github.com/TudorHulban/bCRM/structs"
)

// IAccount Interface provides decoupling of persistence solution for account related operations.
type IAccount interface {
	CreateUser(structs.User) error
	GetUserByUserCode(string) (structs.User, error)
}
