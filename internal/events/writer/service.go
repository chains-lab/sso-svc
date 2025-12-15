package writer

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/chains-lab/sso-svc/internal/events/outbox"
	"github.com/chains-lab/sso-svc/internal/repo"
)

type Service struct {
	repository repository
	addr       string
}

type repository interface {
	CreateOutboxEvent(
		ctx context.Context,
		params repo.CreateOutboxEventParams,
	) (outbox.OutboxEvent, error)
}

func New(addr string, repository repository) *Service {
	return &Service{
		repository: repository,
		addr:       addr,
	}
}

func (s Service) addToOutbox(
	ctx context.Context,
	params contracts.Event,
) error {
	payloadBytes, err := json.Marshal(params.Payload)
	if err != nil {
		return err
	}

	_, err = s.repository.CreateOutboxEvent(
		ctx,
		repo.CreateOutboxEventParams{
			Topic:        params.Topic,
			EventType:    params.EventType,
			EventVersion: params.EventVersion,
			Key:          params.Key,
			Payload:      payloadBytes,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
