package service

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config
	Domain domain.Domain
	Log    *logrus.Logger
}

func NewService(cfg *config.Config, dmn domain.Domain, log *logrus.Logger) (*Service, error) {
	return &Service{
		Config: cfg,
		Log:    log,
		Domain: dmn,
	}, nil
}
