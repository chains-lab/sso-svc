package sql

import (
	"database/sql"

	repositories2 "github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories"
	"github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories/sqlcore"
)

type Repo struct {
	Accounts repositories2.Accounts
	Sessions repositories2.Sessions
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

func NewRepoSQL(url string) (*Repo, error) {
	db, err := NewDBConnection(url)
	if err != nil {
		return nil, err
	}
	queries := sqlcore.New(db)
	return &Repo{
		Accounts: repositories2.NewAccount(queries),
		Sessions: repositories2.NewSession(queries),
	}, nil
}
