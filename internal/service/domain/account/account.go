package account

import (
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx"
	"github.com/sirupsen/logrus"
)

type Account interface {
}

type account struct {
	repo   dbx.Accounts
	Logger *logrus.Logger
}

func NewAccount(accRepo dbx.Accounts, logger *logrus.Logger) Account {
	return &account{
		repo:   accRepo,
		Logger: logger,
	}
}
