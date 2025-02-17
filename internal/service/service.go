package service

import (
	"time"

	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/bin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Service struct {
	Config       *config.Config
	Logger       *logrus.Logger
	TokenManager tokens.TokenManager
	GoogleOAuth  oauth2.Config
	Rabbit       rerabbit.RabbitBroker
	DB           *data.Data
}

func NewService(cfg *config.Config, logger *logrus.Logger) (*Service, error) {
	broker, err := rerabbit.NewBroker(cfg.Rabbit.URL)
	if err != nil {
		return nil, err
	}

	dataBase, err := data.NewDataBase(cfg)
	if err != nil {
		return nil, err
	}

	tokenBin := bin.NewUsersBin(cfg.JWT.Bin.Addr, cfg.JWT.Bin.Password, cfg.JWT.Bin.DB, time.Duration(cfg.JWT.Bin.Lifetime))
	tm := tokens.NewTokenManager(tokenBin, cfg.JWT.AccessToken.SecretKey)

	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)

	return &Service{
		Config:       cfg,
		Logger:       logger,
		TokenManager: tm,
		GoogleOAuth:  googleOAuth,
		Rabbit:       broker,
		DB:           dataBase,
	}, nil
}
