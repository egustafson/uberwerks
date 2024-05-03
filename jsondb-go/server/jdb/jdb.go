package jdb

import (
	"github.com/egustafson/uberwerks/jsondb-go/jsondb"
)

// JDB is the principal interface to direct access of a JSON-DB
type JDB interface {
	List() ([]jsondb.JID, error)
	Put(jo jsondb.JSONObj) (jsondb.JSONObj, error)
	Get(id jsondb.JID) (jo jsondb.JSONObj, ok bool, err error)
	Del(id jsondb.JID) (ok bool, err error)

	Dump() ([]jsondb.JSONObj, error)
	Load(data []map[string]any) error
	Ping() error
}

// Init initializes the JSON-DB.
func Init() error {
	//
	// nothing TODO -- yet
	//
	return nil
}
