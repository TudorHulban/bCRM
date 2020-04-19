package pgstore

import (
	"database/sql"

	"github.com/TudorHulban/bCRM/constants"
	"github.com/TudorHulban/bCRM/pkg/helpers"
	"github.com/TudorHulban/bCRM/structs"
)

// Needs cache for users. Not to go to db for user ID.

var userRights map[int]string

// Userpg type would satisfy IAccount interface.
// Added new type for local decoupling.
type Userpg structs.User

// CreateUser Saves the user variable in the Pg layer. Pointer needed as ID would be read from RDBMS insert.
func (u *Userpg) CreateUser(userData *Userpg, db *sql.DB) error {
	salt := helpers.GenerateRandomString(constants.SaltLength)

	hash, errHash := helpers.HashPassword(userData.LoginPWD, salt)
	if errHash != nil {
		return errHash
	}
	userData.PasswordSALT = salt
	userData.PasswordHASH = hash

	for _, v := range userData.ContactInfo {
		errInsertContact := b.DBConn.Insert(v)
		if errInsertContact != nil {
			return errInsertContact
		}
		pUser.ContactIDs = append(pUser.ContactIDs, v.ID)
	}
	errInsertUser := b.DBConn.Insert(pUser)
	if errInsertUser != nil {
		return errInsertUser
	}

	for _, v := range pUser.ContactInfo {
		v.UserID = pUser.ID
		errUpdateContact := b.DBConn.Update(v)
		if errUpdateContact != nil {
			return errUpdateContact
		}
	}
	return nil
}

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