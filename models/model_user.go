package models

import (
	"context"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo"
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
	}
	c.Logger().Debugf("database is responding.")

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
	u.log.Debugf("user data to insert: %v", u.UserData.UserFormData)

	salt := GenerateRandomString(commons.SaltLength)
	u.UserData.PasswordSALT = salt

	hash, errHash := HashPassword(u.UserData.UserFormData.LoginPWD, u.UserData.PasswordSALT)
	if errHash != nil {
		return errHash
	}
	u.UserData.PasswordHASH = hash

	/*
		for _, v := range userData.ContactInfo {
			errInsertContact := dbConn.Insert(v)
			if errInsertContact != nil {
				return errInsertContact
			}
			userData.ContactIDs = append(userData.ContactIDs, v.ID)
		}
	*/

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	if errInsert := u.db.WithContext(ctx).Insert(&u.UserData); errInsert != nil {
		return errInsert
	}

	/*
		// based on the ID of the row inserted
		for _, v := range userData.ContactInfo {
			v.UserID = userData.ID

			errUpdateContact := dbConn.Update(v)
			if errUpdateContact != nil {
				return errUpdateContact
			}
		}
	*/
	return nil
}

// GetbyID Method based on user ID fetches full user info.
func (u *User) GetbyID(ctx context.Context, timeoutSecs int, userID int64) (UserData, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	result := UserData{}
	errSelect := u.db.WithContext(ctx).Model(&result).Where("id = ?", userID).Select()

	u.log.Debug("fetched:", result)
	return result, errSelect
}

/*
// GetUserByPK Method fetches user info from Pg and returns a user and error.
func (u *User) GetUserByID(userID int64) (*User, error) {

	result := User{ID: pID}

	// verify if requester
	requester, errSelectRequester := getRequesterSecurityGroup(b, 1)
	if errSelectRequester != nil {
		return result, errSelectRequester
	}

	var errSelect error
	switch requester.SecurityGroup {
	case 1:
		{
			errSelect = b.DBConn.Select(&result)
		}
	case 2:
		{
			errSelect = b.DBConn.Model(&result).Where("team_id = ?", requester.TeamID).Where("id = ?", pID).Select()
		}
	}
	if errSelect != nil {
		return result, errSelect
	}
	return result, getContactInfo(b, &result)
}

// GetUserByCode retrieves user given code.
func (u *Userpg) GetUserByCode(pRequesterUserID int64, pCODE string) (Userpg, error) {
	result := User{LoginCODE: pCODE}
	requester, errSelectRequester := getRequesterSecurityGroup(b, 1)
	if errSelectRequester != nil {
		return result, errSelectRequester
	}
	var errSelect error
	switch requester.SecurityGroup {
	case 1:
		{
			errSelect = b.DBConn.Model(&result).Where("login_code = ?", pCODE).Select()
		}
	default:
		{
			errSelect = b.DBConn.Model(&result).Where("team_id = ?", requester.TeamID).Where("login_code = ?", pCODE).Select()
		}
	}
	if errSelect != nil {
		return result, errSelect
	}
	return result, getContactInfo(b, &result)
}

// GetUserByCodeUnauthorized retrieves user given code.
func (u *Userpg) GetUserByCodeUnauthorized(pCODE string) (Userpg, error) {
	result := User{LoginCODE: pCODE}
	errSelect := b.DBConn.Model(&result).Where("login_code = ?", pCODE).Select()

	if errSelect != nil {
		return result, errSelect
	}
	return result, getContactInfo(b, &result)
}

// GetAllUsers retrieves user as per requester security rights.
func (u *Userpg) GetAllUsers(pRequesterUserID int64, pHowMany int) ([]Userpg, error) {
	var result []User
	requester, errSelectRequester := getRequesterSecurityGroup(b, 1)
	if errSelectRequester != nil {
		return result, errSelectRequester
	}
	var errSelect error
	switch requester.SecurityGroup {
	case 1:
		{
			errSelect = b.DBConn.Model(&result).Order("id DESC").Limit(pHowMany).Select()
		}
	case 2:
		{
			errSelect = b.DBConn.Model(&result).Where("team_id = ?", requester.TeamID).Limit(pHowMany).Select()
		}
	}
	return result, errSelect
}

func (u *Userpg) GetMaxIDUsers() (int64, error) {
	var maxID struct {
		Max int64
	}
	_, errQuery := b.DBConn.QueryOne(&maxID, "select max(id) from users")
	return maxID.Max, errQuery
}

func (u *Userpg) UpdateUser(pUser *Userpg) error {
	return b.DBConn.Update(pUser)
}

func getContactInfo(b *Blog, pUser *Userpg) error {
	for _, v := range pUser.ContactIDs {
		co := new(Contact)
		co.ID = v
		errSelectContact := b.DBConn.Select(co)
		if errSelectContact != nil {
			return errSelectContact
		}
		pUser.ContactInfo = append(pUser.ContactInfo, co)
	}
	return nil
}
*/
