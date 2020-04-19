package store

import (
	"github.com/TudorHulban/bCRM/interfaces"
)

var gstore interfaces.IStore

// TheGeneralStore ....
func TheGeneralStore() interfaces.IStore {
	if gstore != nil {
		return gstore
	}
	return nil
}
