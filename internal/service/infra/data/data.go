package data

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/data/repository"
	"github.com/sirupsen/logrus"
)

type Data struct {
	//TODO Переделать чтобы redis db создавался 1 екзепляр и передавался а не создавалось 2 внутрифункций
	Accounts repository.Accounts
	Sessions repository.Sessions
}

func NewData(cfg *config.Config, log *logrus.Logger) (*Data, error) {
	acc, err := repository.NewAccounts(cfg, log)
	if err != nil {
		return nil, err
	}
	sess, err := repository.NewSessions(cfg, log)
	if err != nil {
		return nil, err
	}

	return &Data{
		Accounts: acc,
		Sessions: sess,
	}, nil
}
