package handlers

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App struct {
	Config      *config.Config
	Domain      domain.Domain
	GoogleOAuth oauth2.Config
	Log         *logrus.Logger
}

func NewApp(cfg *config.Config, log *logrus.Logger) (*App, error) {
	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)
	logic, err := domain.NewDomain(cfg, log)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:      cfg,
		Domain:      logic,
		GoogleOAuth: googleOAuth,
		Log:         log,
	}, nil
}
