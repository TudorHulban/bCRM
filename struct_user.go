package main

// User is the representation of the user of the app in the Postgres persistence layer.
// Several methods are defined on this structure in order to satisfy RDBMSUser interface.
// Sorted for maligned.
type User struct {
	ID                  int64 `json:"ID"` // primary key, provided after insert thus pointer needed.
	TeamID              int   // security groups 2, 3 can only see teams tickets
	SecurityGroup       int   `pg:",notnull"` // as per userRights, userRights = map[int]string{1: "admin", 2: "user", 3: "external user"}
	AssignedOpenTickets int   // number of assigned tickets

	PasswordSALT string `json:"-" pg:",notnull` // should not be sent in JSON, exported for ORM
	PasswordHASH string `json:"-" pg:",notnull` // should not be sent in JSON, exported for ORM
	LoginCODE    string `json:"code" pg:",notnull,unique"`
	LoginPWD     string `json:"-" pg:",notnull` // should not be sent in JSON, exported for ORM

	ContactIDs  []int64    // user should accommodate several contacts
	ContactInfo []*Contact `pg:"-"` // when user is retrieved the slice would contain the contacts
}
