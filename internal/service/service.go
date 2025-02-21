package service

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config
	Log    *logrus.Logger
}

func NewService(cfg *config.Config, log *logrus.Logger) (*Service, error) {
	return &Service{
		Config: cfg,
		Log:    log,
	}, nil
}
