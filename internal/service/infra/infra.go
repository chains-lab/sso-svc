package infra

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/jwtmanager"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Accounts repository.Accounts
	Sessions repository.Sessions

	Tokens jwtmanager.JWTManager
}

func NewDataBase(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	acc, err := repository.NewAccounts(cfg, log)
	if err != nil {
		return nil, err
	}
	sess, err := repository.NewSessions(cfg, log)
	if err != nil {
		return nil, err
	}

	jwtManager := jwtmanager.NewJWTManager(cfg)

	return &Infra{
		Accounts: acc,
		Sessions: sess,
		Tokens:   jwtManager,
	}, nil
}
