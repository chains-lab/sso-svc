package transport

import (
	"time"

	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/bin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Transport struct {
	Config       *config.Config
	GoogleOAuth  oauth2.Config
	TokenManager tokens.TokenManager
	Domain       domain.Domain
}

func NewTransport(cfg *config.Config, log *logrus.Logger) (*Transport, error) {
	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)

	tokenBin := bin.NewUsersBin(cfg.JWT.Bin.Addr, cfg.JWT.Bin.Password, cfg.JWT.Bin.DB, time.Duration(cfg.JWT.Bin.Lifetime))
	tm := tokens.NewTokenManager(tokenBin, cfg.JWT.AccessToken.SecretKey)

	bus, err := domain.NewDomain(cfg, log)
	if err != nil {
		return nil, err
	}

	return &Transport{
		Config:       cfg,
		GoogleOAuth:  googleOAuth,
		TokenManager: tm,
		Domain:       bus,
	}, nil
}
