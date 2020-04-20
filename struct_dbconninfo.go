package main

// DBConnInfo Structure holding information for creating RDBMS connection.
type DBConnInfo struct {
	Socket string
	User   string
	Pass   string
	DB     string
}
