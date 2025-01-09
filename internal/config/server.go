package config

import (
	"github.com/cifra-city/mailman"
	"github.com/cifra-city/sso-oauth/internal/data/sql"
	"github.com/cifra-city/tokens"
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
	Mailman      *mailman.Mailman
	TokenManager tokens.TokenManager
	Broker       *Broker
	GoogleOAuth  oauth2.Config
}

func NewServer(cfg *Config) (*Server, error) {
	logger := SetupLogger(cfg.Logging.Level, cfg.Logging.Format)
	mail := mailman.NewMailman(cfg.Email.SmtpPort, cfg.Email.SmtpHost, cfg.Email.Address, cfg.Email.Password)
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
		Mailman:      mail,
		TokenManager: tokenManager,
		Broker:       broker,
		GoogleOAuth:  googleOAuth,
	}, nil
}
