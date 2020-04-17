package persistance

import (
	"github.com/TudorHulban/bCRM/pkg/structs"
)

// IAccount Interface provides decoupling of persistance solution for account related operations.
type IAccount interface {
	GetUserByUserCode(code string) (structs.User, error)
}
