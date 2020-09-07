package sql

import (
	_ "github.com/denisenkom/go-mssqldb"
	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
)

type Store struct {
	DB *dbx.DB
}

func New(config *config.Database) *Store {
	db, err := dbx.MustOpen("mssql", config.DSN())
	if err != nil {
		panic(err)
	}

	return &Store{DB: db}
}

func (s *Store) Shutdown() error {
	return s.DB.Close()
}
