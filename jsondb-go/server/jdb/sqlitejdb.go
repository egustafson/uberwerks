package jdb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// sqliteJDB is an SQLite implementation of a JDB
type sqliteJDB struct {
	SqlJDB
	dsn string
}

func InitSqliteJDB(dsn string) (JDB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		db.Close()
		return nil, err
	}
	return &sqliteJDB{
		SqlJDB: SqlJDB{DB: db},
		dsn:    dsn,
	}, nil
}
