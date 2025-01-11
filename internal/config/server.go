package config

import (
	"github.com/recovery-flow/sso-oauth/internal/data/sql"
	"github.com/recovery-flow/tokens"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	SERVER = "service"
)

type Server struct {
	Config       *Config
	SqlDB        *sql.Repo
	Logger       *logrus.Logger
	TokenManager tokens.TokenManager
	Broker       *Broker
	GoogleOAuth  oauth2.Config
}

func NewServer(cfg *Config) (*Server, error) {
	logger := SetupLogger(cfg.Logging.Level, cfg.Logging.Format)
	queries, err := sql.NewRepoSQL(cfg.Database.URL)
	tokenManager := tokens.NewTokenManager(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logger, cfg.JWT.AccessToken.TokenLifetime)
	broker, err := NewBroker(cfg.Rabbit.URL, cfg.Rabbit.Exchange)
	if err != nil {
		return nil, err
	}
	googleOAuth := InitGoogleOAuth(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)

	return &Server{
		Config:       cfg,
		SqlDB:        queries,
		Logger:       logger,
		TokenManager: tokenManager,
		Broker:       broker,
		GoogleOAuth:  googleOAuth,
	}, nil
}
