package models

// TeamMembers Structure holding team members.
type TeamMembers struct {
	tableName struct{} `pg:"teamsmembers"`
	ID        int64
	TeamID    int64 `validate:"required" pg:",notnull,unique"`
	UserID    int64 `validate:"required" pg:",notnull,unique"`
	Joined    int64 // unix time seconds when user joined team
	JoinedBy  int64 // user ID that added to team
	Left      int64 // unix time seconds when user left team
	LeftBy    int64 // user ID that eliberated user from team
}

type Team struct {
	TeamFormData
	tools
}
