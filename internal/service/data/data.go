package data

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data/repositories"
)

type Data struct {
	Accounts repositories.Accounts
	Sessions repositories.Sessions
}

func NewDataBase(cfg *config.Config) (*Data, error) {
	acc, err := repositories.NewAccounts(cfg)
	if err != nil {
		return nil, err
	}
	sess, err := repositories.NewSessions(cfg)
	if err != nil {
		return nil, err
	}

	return &Data{
		Accounts: acc,
		Sessions: sess,
	}, nil
}
