package infra

import (
	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/data"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/jwtmanager"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Tokens jwtmanager.JWTManager
	Rabbit rerabbit.RabbitBroker
	Data   *data.Data
}

func NewInfra(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	eve, err := rerabbit.NewBroker(cfg.Rabbit.URL)
	if err != nil {
		return nil, err
	}
	NewData, err := data.NewData(cfg, log)
	if err != nil {
		return nil, err
	}
	jwtManager := jwtmanager.NewJWTManager(cfg)
	return &Infra{
		Data:   NewData,
		Tokens: jwtManager,
		Rabbit: eve,
	}, nil
}
