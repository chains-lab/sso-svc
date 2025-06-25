package app

import (
	"github.com/chains-lab/chains-auth/internal/app/domain"
	"github.com/chains-lab/chains-auth/internal/utils/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	sessions sessionsDomain
	users    usersDomain
	log      *logrus.Entry
}

func NewApp(cfg config.Config, log *logrus.Logger) (App, error) {
	sessions, err := domain.NewSession(cfg, log)
	if err != nil {
		log.WithError(err).Error("failed to create sessions domain")
		return App{}, err
	}

	users, err := domain.NewUser(cfg, log)
	if err != nil {
		log.WithError(err).Error("failed to create users domain")
		return App{}, err
	}

	return App{
		sessions: sessions,
		users:    users,
		log:      log.WithField("component", "app"),
	}, nil
}
