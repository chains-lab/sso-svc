package infra

import (
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/service/infra/data"
	"github.com/hs-zavet/sso-oauth/internal/service/infra/events/producer"
	"github.com/hs-zavet/sso-oauth/internal/service/infra/jwtmanager"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Tokens jwtmanager.JWTManager
	Kafka  producer.Producer
	Data   *data.Data
}

func NewInfra(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	jwtManager := jwtmanager.NewJWTManager(cfg)
	eve := producer.NewProducer(cfg)
	NewData, err := data.NewData(cfg, log)
	if err != nil {
		return nil, err
	}
	return &Infra{
		Tokens: jwtManager,
		Kafka:  eve,
		Data:   NewData,
	}, nil
}
