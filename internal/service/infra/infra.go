package infra

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	sqldb2 "github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb"
)

type Data struct {
	Accounts sqldb2.Accounts
	Sessions sqldb2.Sessions
}

func NewDataBase(cfg *config.Config) (*Data, error) {
	acc, err := sqldb2.NewAccounts(cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}
	sess, err := sqldb2.NewSessions(cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}

	return &Data{
		Accounts: acc,
		Sessions: sess,
	}, nil
}
