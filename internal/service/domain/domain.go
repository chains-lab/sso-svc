package domain

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/entities"
	"github.com/sirupsen/logrus"
)

type Domain struct {
	Log     *logrus.Logger
	Session entities.Session
	Account entities.Account
}

func NewDomain(cfg *config.Config, logger *logrus.Logger) (*Domain, error) {
	acc := entities.NewAccount(DB.Accounts, logger)
	ses := entities.NewSession(DB.Sessions, logger)

	return &Domain{
		Log:     logger,
		Session: ses,
		Account: acc,
	}, err
}
