package service

import (
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/data"
	"github.com/recovery-flow/tokens"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Service struct {
	Config       *config.Config
	Logger       *logrus.Logger
	TokenManager *tokens.TokenManager
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

	tm := tokens.NewTokenManager(cfg.Database.Redis.Addr, cfg.Database.Redis.Password, cfg.Database.Redis.DB, logger, cfg.JWT.AccessToken.TokenLifetime)
	googleOAuth := config.InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)

	return &Service{
		Config:       cfg,
		Logger:       logger,
		TokenManager: &tm,
		GoogleOAuth:  googleOAuth,
		Rabbit:       broker,
		DB:           dataBase,
	}, nil
}
