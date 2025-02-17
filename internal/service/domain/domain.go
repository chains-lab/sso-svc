package domain

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/account"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/session"
	"github.com/sirupsen/logrus"
)

type Domain struct {
	Log  *logrus.Logger
	Data *data.Data

	Session session.Session
	Account account.Account
}

func NewDomain(cfg *config.Config, logger *logrus.Logger) (*Domain, error) {
	DB, err := data.NewDataBase(cfg)
	if err != nil {
		return nil, err
	}

	acc := account.NewAccount(DB.Accounts, logger)
	ses := session.NewSession(DB.Sessions, logger)

	return &Domain{
		Log:     logger,
		Data:    DB,
		Session: ses,
		Account: acc,
	}, err
}
