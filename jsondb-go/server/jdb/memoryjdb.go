package jdb

import (
	"github.com/egustafson/uberwerks/jsondb-go/jsondb"
	"github.com/google/uuid"
)

// memoryJDB is an in-memory implementation of a JDB
type memoryJDB struct {
	db map[jsondb.JID]jsondb.JSONObj
}

// static check: a memoryJDB ptr isA JDB interface
var _ JDB = (*memoryJDB)(nil)

func NewMemoryJDB() JDB {
	return &memoryJDB{
		db: make(map[jsondb.JID]jsondb.JSONObj),
	}
}

func InitMemoryJDB(data []map[string]any) JDB {
	jdb := NewMemoryJDB()
	_ = jdb.Load(data)
	return jdb
}

func (jdb *memoryJDB) Ping() error {
	return nil
}

func (jdb *memoryJDB) List() ([]jsondb.JID, error) {
	l := make([]jsondb.JID, 0, len(jdb.db))
	for k := range jdb.db {
		l = append(l, jsondb.JID(k))
	}
	return l, nil
}

func (jdb *memoryJDB) Put(jo jsondb.JSONObj) (jsondb.JSONObj, error) {
	jid, ok := jo[jsondb.IDKey].(string)
	if !ok {
		jid = uuid.NewString()
		jo[jsondb.IDKey] = jid
	}
	jdb.db[jsondb.JID(jid)] = jo
	return jo, nil
}

func (jdb *memoryJDB) Get(id jsondb.JID) (jo jsondb.JSONObj, ok bool, err error) {
	jo, ok = jdb.db[id]
	return
}

func (jdb *memoryJDB) Del(id jsondb.JID) (ok bool, err error) {
	if _, ok = jdb.db[id]; ok {
		delete(jdb.db, id)
	}
	return
}

func (jdb *memoryJDB) Load(data []map[string]any) error {
	for _, obj := range data {
		jobj := make(map[jsondb.JID]any)
		for k, v := range obj {
			jobj[jsondb.JID(k)] = v
		}
		jdb.Put(jobj)
	}
	return nil
}

func (jdb *memoryJDB) Dump() ([]jsondb.JSONObj, error) {
	objs := make([]jsondb.JSONObj, 0, len(jdb.db))
	for _, v := range jdb.db {
		objs = append(objs, v)
	}
	return objs, nil
}
