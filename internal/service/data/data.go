package data

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx"
)

type Data struct {
	Accounts dbx.Accounts
	Sessions dbx.Sessions
}

func NewDataBase(cfg *config.Config) (*Data, error) {
	acc, err := dbx.NewAccounts(cfg)
	if err != nil {
		return nil, err
	}
	sess, err := dbx.NewSessions(cfg)
	if err != nil {
		return nil, err
	}

	return &Data{
		Accounts: acc,
		Sessions: sess,
	}, nil
}
