package handlers

import (
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"golang.org/x/oauth2"
)

type Handler struct {
	app    App
	log    logium.Logger
	cfg    config.Config
	google oauth2.Config
}

func NewHandler(cfg config.Config, log logium.Logger, a *app.App) Handler {
	return Handler{
		app:    a,
		log:    log,
		cfg:    cfg,
		google: cfg.GoogleOAuth(),
	}
}
