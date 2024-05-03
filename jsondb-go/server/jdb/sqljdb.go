package jdb

import (
	"database/sql"
	"encoding/json"

	"github.com/egustafson/uberwerks/jsondb-go/jsondb"
)

const (
	create string = `
      CREATE TABLE IF NOT EXISTS jsondb (
	    id STRING NOT NULL PRIMARY KEY,
		data TEXT
  	  );`

	select_json  string = `SELECT data FROM jsondb;`
	select_ids   string = `SELECT id FROM jsondb;`
	select_by_id string = `SELECT data FROM jsondb WHERE id=?`
	insert       string = `INSERT INTO jsondb VALUES(?,?);`
	delete_by_id string = `DELETE FROM jsondb WHERE id=?`
)

type SqlJDB struct {
	DB *sql.DB
}

// static check: a sqlJDB ptr isA JDB interface
var _ JDB = (*SqlJDB)(nil)

func NewSqlJDB(db *sql.DB) JDB {
	return &SqlJDB{DB: db}
}

func (jdb *SqlJDB) Ping() error {
	return jdb.DB.Ping()
}

func (jdb *SqlJDB) List() ([]jsondb.JID, error) {
	rows, err := jdb.DB.Query(select_ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]jsondb.JID, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		list = append(list, jsondb.JID(id))
	}
	return list, nil
}

func (jdb *SqlJDB) Put(jo jsondb.JSONObj) (jsondb.JSONObj, error) {
	id := jsondb.Identify(jo)
	jsonb, err := json.Marshal(jo)
	if err != nil {
		return nil, err
	}
	_, err = jdb.DB.Exec(insert, string(id), string(jsonb))
	if err != nil {
		return nil, err
	}
	return jo, nil
}

func (jdb *SqlJDB) Get(id jsondb.JID) (jo jsondb.JSONObj, ok bool, err error) {
	var jsonb []byte
	err = jdb.DB.QueryRow(select_by_id, string(id)).Scan(&jsonb)
	switch {
	case err == sql.ErrNoRows:
		ok = false
		return nil, false, nil
	case err != nil:
		return nil, false, err
	}
	jo = make(jsondb.JSONObj)
	if err := json.Unmarshal(jsonb, &jo); err != nil {
		return nil, false, err
	}
	return jo, true, nil
}

func (jdb *SqlJDB) Del(id jsondb.JID) (ok bool, err error) {
	var result sql.Result
	if result, err = jdb.DB.Exec(delete_by_id, string(id)); err != nil {
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	return (rowsAffected > 0), nil
}

func (jdb *SqlJDB) Dump() ([]jsondb.JSONObj, error) {
	rows, err := jdb.DB.Query(select_json)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]jsondb.JSONObj, 0)
	for rows.Next() {
		jsonb := make([]byte, 0)
		if err = rows.Scan(&jsonb); err != nil {
			return nil, err
		}
		obj := make(jsondb.JSONObj)
		if err = json.Unmarshal(jsonb, &obj); err != nil {
			return nil, err
		}
		list = append(list, obj)
	}
	return list, nil
}

func (jdb *SqlJDB) Load(data []map[string]any) error {
	for _, robj := range data {
		obj := make(jsondb.JSONObj)
		for k, v := range robj {
			obj[jsondb.JID(k)] = v
		}
		if _, err := jdb.Put(obj); err != nil {
			return err
		}
	}
	return nil
}
