package handlers

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service"
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

func NewHandlers(svc *service.Service) (*Handlers, error) {
	googleOAuth := config.InitGoogleOAuth(svc.Config.OAuth.Google.ClientID, svc.Config.OAuth.Google.ClientSecret, svc.Config.OAuth.Google.RedirectURL)

	return &Handlers{
		Config:      svc.Config,
		Domain:      svc.Domain,
		GoogleOAuth: googleOAuth,
		Log:         svc.Log,
	}, nil
}
