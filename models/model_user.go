package models

import (
	"context"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Needs cache for users. Not to go to db for user ID.

// Contact Is used when defining a app user. The app user could have more than one contact.
type Contact struct {
	ID             int64
	UserID         int64  `pg:"userid"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lasttname"`
	OfficePhoneNo  string
	MobilePhoneNo  string
	CompanyEmail   string
	WorkEmail      string
	AddressHQ      string
	AddressOffice  string
	AddressBilling string
}

// UserFormData Structure holds information necessary for creating a user and coming from frontend.
type UserFormData struct {
	TeamID        int    `validate:"required"`               // triggers an entry in teams data table.
	SecurityGroup int    `pg:",notnull" validate:"required"` // as per userRights, userRights = map[int]string{1: "admin", 2: "user", 3: "external user"}
	LoginCODE     string `validate:"required" json:"code" pg:",notnull,unique" `
	LoginPWD      string `validate:"required" json:"-" pg:",notnull ` // should not be sent in JSON, exported for ORM, to be taken out as hash is enough
}

// UserData Structure holds the actual user persisted user information.
type UserData struct {
	tableName struct{} `pg:"users"`
	ID        int64    `json:"ID" valid:"-"` // primary key, provided after insert thus pointer needed.
	UserFormData

	AssignedOpenTickets int    `valid:"-"`                                             // number of assigned tickets
	PasswordSALT        string `valid:"type(string), optional" json:"-" pg:",notnull ` // should not be sent in JSON, exported for ORM
	PasswordHASH        string `valid:"type(string)" json:"-" pg:",notnull `           // should not be sent in JSON, exported for ORM

	//ContactIDs  []int64    `valid:"type(string), optional"` // user should accommodate several contacts
	//ContactInfo []*Contact `pg:"-" valid:"-"`               // when user is retrieved the slice would contain the contacts
}

// User is the representation of the user of the app in the Postgres persistence layer.
// Several methods are defined on this structure in order to satisfy RDBMSUser interface.
type User struct {
	UserData
	tools
}

var userRights map[int]string

// NewUser Constructor for when interacting with the user model.
// Use validation for inserts or updates. No validation for selects.
func NewUser(c echo.Context, f UserFormData, noValidation bool) (*User, error) {
	// validate data
	if !noValidation {
		errValid := isValidStruct(f, c.Logger())
		if errValid != nil {
			return nil, errValid
		}
	}

	// check db connection. debug level = 1
	if c.Logger().Level() == 1 {
		errQuery := commons.CheckPgDB(c.Logger())
		if errQuery != nil {
			return nil, errQuery
		}
		c.Logger().Debugf("database is responding.")
	}

	return &User{
		UserData: UserData{UserFormData: f},
		tools: tools{
			log: c.Logger(),
			db:  commons.DB(),
		},
	}, nil
}

// CreateUser Saves the user variable in the Pg layer. Pointer needed as ID would be read from RDBMS insert.
func (u *User) Insert(ctx context.Context, timeoutSecs int) error {
	salt := GenerateRandomString(commons.SaltLength)
	u.UserData.PasswordSALT = salt

	hash, errHash := HashPassword(u.UserData.UserFormData.LoginPWD, u.UserData.PasswordSALT)
	if errHash != nil {
		return errHash
	}
	u.UserData.PasswordHASH = hash

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	u.log.Debugf("user data to insert: %v", u.UserData.UserFormData)

	u.tools.log.SetLevel(log.DEBUG)
	// check db connection. debug level = 1
	if u.tools.log.Level() == 1 {
		errQuery := commons.CheckPgDB(u.tools.log)
		if errQuery != nil {
			return errQuery
		}
		u.tools.log.Debugf("database is responding.")
	}
	u.tools.log.Info("log level: ", u.log.Level())

	if errInsert := u.db.WithContext(ctx).Insert(&u.UserData); errInsert != nil {
		return errInsert
	}
	return nil
}
