package handlers

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Handlers struct {
	Config      *config.Config
	Domain      domain.Domain
	GoogleOAuth oauth2.Config
	Log         *logrus.Logger
}

func NewHandlers(cfg *config.Config, log *logrus.Logger) (*Handlers, error) {
	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)
	logic, err := domain.NewDomain(cfg, log)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		Config:      cfg,
		Domain:      logic,
		GoogleOAuth: googleOAuth,
		Log:         log,
	}, nil
}
