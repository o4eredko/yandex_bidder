package sql

import (
	_ "github.com/denisenkom/go-mssqldb"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type (
	Store struct {
		DB *dbx.DB
	}
)

func New(connString string) *Store {
	db, err := dbx.MustOpen("mssql", connString)
	if err != nil {
		panic(err)
	}

	return &Store{DB: db}
}

func (s *Store) Shutdown() error {
	return s.DB.Close()
}
