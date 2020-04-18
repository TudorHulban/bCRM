package structs

// Team Structure holding team fields. Every application user belongs to a team.
type Team struct {
	ID              int64
	CODE            string `pg:",notnull,unique"`
	Name            string
	Description     string
	ManagerID       int64 `pg:"managerid"`
	AssignedTickets int   // number of assigned tickets
}
