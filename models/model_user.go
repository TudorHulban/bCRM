package models

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"

	"github.com/TudorHulban/bCRM/pkg/commons"
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
	TeamID        int    `pg:",notnull" valid:"type(int)"` // security groups 2, 3 can only see teams tickets
	SecurityGroup int    `pg:",notnull" valid:"type(int)"` // as per userRights, userRights = map[int]string{1: "admin", 2: "user", 3: "external user"}
	LoginCODE     string `valid:"type(string)" json:"code" pg:",notnull,unique" `
	LoginPWD      string `valid:"type(string)" json:"-" pg:",notnull ` // should not be sent in JSON, exported for ORM, to be taken out as hash is enough
}

// User is the representation of the user of the app in the Postgres persistence layer.
// Several methods are defined on this structure in order to satisfy RDBMSUser interface.
// Sorted for maligned.
type User struct {
	ID int64 `json:"ID" valid:"-"` // primary key, provided after insert thus pointer needed.
	UserFormData

	AssignedOpenTickets int    `valid:"-"`                                             // number of assigned tickets
	PasswordSALT        string `valid:"type(string), optional" json:"-" pg:",notnull ` // should not be sent in JSON, exported for ORM
	PasswordHASH        string `valid:"type(string)" json:"-" pg:",notnull `           // should not be sent in JSON, exported for ORM

	//ContactIDs  []int64    `valid:"type(string), optional"` // user should accommodate several contacts
	//ContactInfo []*Contact `pg:"-" valid:"-"`               // when user is retrieved the slice would contain the contacts

	valid *validator.Validate
	log   echo.Logger
	db    *pg.DB
}

var userRights map[int]string

func NewUser(c echo.Context, db *pg.DB, f UserFormData) *User {
	return &User{
		UserFormData: f,
		valid:        validator.New(),
		log:          c.Logger(),
		db:           db,
	}
}

// CreateUser Saves the user variable in the Pg layer. Pointer needed as ID would be read from RDBMS insert.
func (u *User) Insert() error {
	u.log.Debugf("user data to insert: ", u.UserFormData)

	errValid := u.valid.Struct(u.UserFormData)
	if errValid != nil {
		u.log.Debugf("validation error:", errValid)
		return errValid
	}
	u.log.Debugf("structure is valid.")

	salt := GenerateRandomString(commons.SaltLength)
	u.PasswordSALT = salt

	hash, errHash := HashPassword(u.UserFormData.LoginPWD, u.PasswordSALT)
	if errHash != nil {
		return errHash
	}
	u.PasswordHASH = hash

	/*
		for _, v := range userData.ContactInfo {
			errInsertContact := dbConn.Insert(v)
			if errInsertContact != nil {
				return errInsertContact
			}
			userData.ContactIDs = append(userData.ContactIDs, v.ID)
		}
	*/

	if errInsertUser := u.db.Insert(u); errInsertUser != nil {
		return errInsertUser
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

/*

// GetUserByPK fetches user info from Pg and returns a user type.
func (u *Userpg) GetUserByPK(pID int64) (Userpg, error) {
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
