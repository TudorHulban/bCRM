package main

// httpError Helper structure for pushing error into JSON.
type httpError struct {
	TheError string `json:"error"` // exported for JSON pick up
}
