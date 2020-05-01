package models

import (
	"context"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// TeamMembers Structure holding team members.
type TeamMembersData struct {
	tableName struct{} `pg:"teamsmembers"`
	ID        int64
	TeamID    int64 `validate:"required" pg:",notnull,unique"`
	UserID    int64 `validate:"required" pg:",notnull,unique"`
	Joined    int64 `pg:",notnull"` // unix time seconds when user joined team
	JoinedBy  int64 `pg:",notnull"` // user ID that added to team
	Left      int64 `pg:""`         // unix time seconds when user left team
	LeftBy    int64 `pg:""`         // user ID that eliberated user from team
}

type TeamMembers struct {
	TeamMembersData
	tools
}

// newTeamMember Only for models package use.
func newTeamMember(l echo.Logger, d TeamMembersData, noValidation bool) (*TeamMembers, error) {
	// validate data
	if !noValidation {
		if errValid := isValidStruct(d, l); errValid != nil {
			return nil, errValid
		}
	}

	result := TeamMembers{
		TeamMembersData: d,
		tools: tools{
			log: l,
			db:  commons.DB(),
		},
	}
	result.tools.log.SetLevel(log.DEBUG)
	return &result, nil
}

func (t *TeamMembers) insert(ctx context.Context, timeoutSecs int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	if errInsert := t.db.WithContext(ctx).Insert(&t.TeamMembersData); errInsert != nil {
		return errInsert
	}
	return nil
}

func (t *TeamMembers) getIDsforUserID(ctx context.Context, timeoutSecs int, userID int64) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	var teams []TeamMembersData
	if errSelect := t.db.WithContext(ctx).Model(&teams).Where("userid = ?", userID).Select(); errSelect != nil {
		return nil, errSelect
	}

	var result []int64
	for _, v := range teams {
		result = append(result, v.TeamID)
	}
	return result, nil
}
