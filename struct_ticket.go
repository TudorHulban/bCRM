package main

// TicketType abstracts a ticket type. Ticket contains events.
type TicketType struct {
	ID          int
	TypeCODE    string
	SLAID       int // service level agreement per type of ticket and customer
	Description string
}

// Ticket concentrates events.
type Ticket struct {
	ID              int64
	TypeCODE        int   // have a mapping of types
	OpenedByUserID  int64 `pg:"userid"`
	OpenedByTeamID  int   `pg:"teamid"`
	Opened          int64 `pg:"openednano"`
	Closed          int64 `pg:"closednano"`
	CurrentStatus   int   // ticket status ID
	CurrentPriority int
	CurrentUserID   int64 `pg:"assignedid"`
	Title           string
	Description     string
	Events          []Event `pg:"-"`
}

// TicketStatus Structure represents the definition of the ticket status.
// To support different flows an option of next tickets statuses is added.
type TicketStatus struct {
	ID                     int
	TicketTypeID           int
	StatusCODE             string
	StatusDescription      string
	NOTAllowedNextStatusID []int64 `pg:"notnextstatusid"`
}

// TicketMovement Structure holding audit information for each movement of tickets.
type TicketMovement struct {
	ID                      int64 // little bit of memory padding
	PreviousStatus          int   // ticket status ID
	ChangedToStatus         int   // ticket status ID
	PreviousPriority        int
	ChangedToPriority       int
	TicketID                int64
	Timestamp               int64
	PreviousAssignedUserID  int64
	ChangedToAssignedUserID int64
}
