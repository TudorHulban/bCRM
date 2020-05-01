package models

import (
	"context"
	"errors"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
)

// getbyID Method based on user ID fetches full user info.
func (u *User) getbyID(ctx context.Context, timeoutSecs int, userID int64) (UserData, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	result := UserData{}
	errSelect := u.db.WithContext(ctx).Model(&result).Where("id = ?", userID).Select()

	u.log.Debug("fetched:", result)
	return result, errSelect
}

func (u *User) getTeams(ctx context.Context, timeoutSecs int, userID int64) ([]int64, error) {
	data := TeamMembersData{
		UserID: userID,
	}
	t, errCo := newTeamMember(u.tools.log, data, true)
	if errCo != nil {
		return nil, errCo
	}
	return t.getIDsforUserID(ctx, timeoutSecs, userID)
}

// GetbyID Method based on user ID and requester ID fetches full user info.
func (u *User) GetbyID(ctx context.Context, timeoutSecs int, userID, requesterID int64) (UserData, error) {
	// get requester info to understand if info should be provided
	requester, errReq := u.getbyID(ctx, timeoutSecs, requesterID)
	if errReq != nil {
		return UserData{}, errReq
	}

	user, errUser := u.getbyID(ctx, timeoutSecs, userID)
	if errUser != nil {
		return UserData{}, errUser
	}

	// check security
	// pass - if requester is app admin
	if requester.SecurityGroup == commons.SecuAppAdmin {
		return user, nil
	}
	// pass - if requester is group admin and in the same group
	if requester.SecurityGroup == commons.SecuGroupAdmin && requester.AppGroup == user.AppGroup {
		return user, nil
	}
	// pass - if requester is team admin and in the same team. to consider the user could be in more than one team.

	return UserData{}, errors.New("data could not be provided due to security issues")
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
