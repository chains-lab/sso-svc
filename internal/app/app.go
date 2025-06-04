package app

import (
	"github.com/chains-lab/chains-auth/internal/app/domain"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	sessions domain.Sessions
	accounts domain.Accounts
	log      *logrus.Entry
}

func NewApp(cfg config.Config, log *logrus.Logger) (App, error) {
	sessions, err := domain.NewSession(cfg, log)
	if err != nil {
		log.WithError(err).Error("failed to create sessions domain")
		return App{}, err
	}

	accounts, err := domain.NewAccount(cfg, log)
	if err != nil {
		log.WithError(err).Error("failed to create accounts domain")
		return App{}, err
	}

	return App{
		sessions: sessions,
		accounts: accounts,
		log:      log.WithField("component", "app"),
	}, nil
}
