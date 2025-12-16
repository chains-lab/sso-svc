package repo

import (
	"database/sql"

	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
)

type Repository struct {
	sql sqlDB
}

type sqlDB struct {
	accounts  pgdb.AccountsQ
	emails    pgdb.AccountEmailsQ
	passwords pgdb.AccountPasswordsQ
	outbox    pgdb.OutboxEventsQ
	sessions  pgdb.SessionsQ
}

func New(db *sql.DB) *Repository {
	return &Repository{
		sql: sqlDB{
			accounts:  pgdb.NewAccounts(db),
			outbox:    pgdb.NewOutboxEvents(db),
			sessions:  pgdb.NewSessions(db),
			emails:    pgdb.NewAccountEmails(db),
			passwords: pgdb.NewAccountPasswords(db),
		},
	}
}
