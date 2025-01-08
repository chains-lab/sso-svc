package db

import (
	"database/sql"

	"github.com/cifra-city/sso-oauth/internal/data/db/sqlcore"
)

type Databaser struct {
	Accounts Accounts
	Sessions Sessions

	Transaction
}

func NewDBConnection(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDatabaser(url string) (*Databaser, error) {
	db, err := NewDBConnection(url)
	if err != nil {
		return nil, err
	}
	queries := sqlcore.New(db)
	return &Databaser{
		Accounts:    NewAccount(queries),
		Sessions:    NewSession(queries),
		Transaction: NewTransaction(queries),
	}, nil
}
