package controller

import (
	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/domain"
	"golang.org/x/oauth2"
)

type Service struct {
	google oauth2.Config
	app    domain.Core
	log    logium.Logger
}

func NewService(log logium.Logger, google oauth2.Config, app domain.Core) *Service {
	return &Service{
		log:    log,
		google: google,
		app:    app,
	}
}
