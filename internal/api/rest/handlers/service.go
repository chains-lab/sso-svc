package handlers

import (
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"golang.org/x/oauth2"
)

type Service struct {
	app    *app.App
	log    logium.Logger
	cfg    config.Config
	google oauth2.Config
}

func NewService(cfg config.Config, log logium.Logger, a *app.App) Service {
	return Service{
		app:    a,
		log:    log,
		cfg:    cfg,
		google: cfg.GoogleOAuth(),
	}
}

func (s Service) Log(r *http.Request) logium.Logger {
	return s.log
}
