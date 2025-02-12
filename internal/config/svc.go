package config

import (
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/data"
	"github.com/recovery-flow/tokens"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	SERVICE = "service"
)

type Service struct {
	Config       *Config
	DB           *data.Data
	Logger       *logrus.Logger
	TokenManager tokens.TokenManager
	Rabbit       rerabbit.RabbitBroker
	GoogleOAuth  oauth2.Config
}

func NewService(cfg *Config) (*Service, error) {
	logger := SetupLogger(cfg.Logging.Level, cfg.Logging.Format)
	tokenManager := tokens.NewTokenManager(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logger, cfg.JWT.AccessToken.TokenLifetime)
	googleOAuth := InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)
	dataBase, err := data.NewDataBase(data.Config{
		SqlUrl:        cfg.Database.URL,
		RedisAddr:     cfg.Redis.Addr,
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
	})
	if err != nil {
		return nil, err
	}
	broker, err := rerabbit.NewBroker(cfg.Rabbit.URL)
	if err != nil {
		return nil, err
	}

	return &Service{
		Config:       cfg,
		DB:           dataBase,
		Logger:       logger,
		TokenManager: tokenManager,
		Rabbit:       broker,
		GoogleOAuth:  googleOAuth,
	}, nil
}
