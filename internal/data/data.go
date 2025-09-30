package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/sso-svc/internal/data/pgdb"
)

type Database struct {
	sql SqlDB
}

func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.sessions.New().Transaction(ctx, fn)
}

type SqlDB struct {
	users    pgdb.UsersQ
	sessions pgdb.SessionsQ
}

func (d *SqlDB) Users() pgdb.UsersQ {
	return d.users
}

func (d *SqlDB) Sessions() pgdb.SessionsQ {
	return d.sessions
}

func NewDatabase(db *sql.DB) *Database {
	usersSql := pgdb.NewUsers(db)
	sessionsSql := pgdb.NewSessions(db)

	return &Database{
		sql: SqlDB{
			users:    usersSql,
			sessions: sessionsSql,
		},
	}
}
