package handlers

import (
	"time"

	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/bin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Handler struct {
	Domain       *domain.Domain
	GoogleOAuth  oauth2.Config
	TokenManager tokens.TokenManager
	Log          *logrus.Logger
}

func NewHandler(cfg *config.Config, log *logrus.Logger) (*Handler, error) {
	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)
	logic, err := domain.NewDomain(cfg, log)

	tokenBin := bin.NewUsersBin(cfg.JWT.Bin.Addr, cfg.JWT.Bin.Password, cfg.JWT.Bin.DB, time.Duration(cfg.JWT.Bin.Lifetime))
	tm := tokens.NewTokenManager(tokenBin, cfg.JWT.AccessToken.SecretKey)
	if err != nil {
		return nil, err
	}
	return &Handler{
		Domain:       logic,
		GoogleOAuth:  googleOAuth,
		TokenManager: tm,
		Log:          log,
	}, nil
}
