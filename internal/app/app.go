package app

import (
	"github.com/chains-lab/sso-svc/internal/app/entities"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/sirupsen/logrus"
)

type App struct {
	sessions entities.Sessions
	users    entities.Users
	log      *logrus.Entry
}

func NewApp(cfg config.Config, log logger.Logger) (App, error) {
	sessions, err := entities.NewSession(cfg, log)
	if err != nil {
		log.WithError(err).Error("failed to create sessions domain")
		return App{}, err
	}

	users, err := entities.NewUser(cfg, log)
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
