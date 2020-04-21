package main

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/gommon/log"
)

// Needs cache for users. Not to go to db for user ID.

var userRights map[int]string

// CreateUser Saves the user variable in the Pg layer. Pointer needed as ID would be read from RDBMS insert.
func CreateUser(userData *User) error {
	log.Print("user data to insert: ", userData)

	salt := GenerateRandomString(SaltLength)

	hash, errHash := HashPassword(userData.LoginPWD, salt)
	if errHash != nil {
		return errHash
	}
	userData.PasswordSALT = salt
	userData.PasswordHASH = hash

	isValid, errValid := valid.ValidateStruct(userData)
	if errValid != nil {
		log.Print("validation error:", errValid, isValid)
		return errValid
	}
	log.Print("structure valid: ", isValid)

	for _, v := range userData.ContactInfo {
		errInsertContact := dbConn.Insert(v)
		if errInsertContact != nil {
			return errInsertContact
		}
		userData.ContactIDs = append(userData.ContactIDs, v.ID)
	}

	errInsertUser := dbConn.Insert(userData)
	if errInsertUser != nil {
		return errInsertUser
	}

	// based on the ID of the row inserted
	for _, v := range userData.ContactInfo {
		v.UserID = userData.ID

		errUpdateContact := dbConn.Update(v)
		if errUpdateContact != nil {
			return errUpdateContact
		}
	}
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
