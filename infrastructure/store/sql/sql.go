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
	db, err := dbx.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}

	return &Store{DB: db}
}

func (s *Store) Shutdown() error {
	if err := s.DB.Close(); err != nil {
		return err
	}

	return nil
}
