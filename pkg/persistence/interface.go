package persistence

import (
	"github.com/TudorHulban/bCRM/pkg/structs"
)

// IAccount Interface provides decoupling of persistence solution for account related operations.
type IAccount interface {
	GetUserByUserCode(code string) (structs.User, error)
}
