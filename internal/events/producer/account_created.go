package producer

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/kafkakit/box"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const AccountCreatedEvent = "account.created"

type AccountCreatedPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email,omitempty"`
}

func (s Service) WriteAccountCreated(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	payload, err := json.Marshal(AccountCreatedPayload{
		Account: account,
		Email:   email,
	})
	if err != nil {
		return err
	}

	eventID := uuid.New()

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.AccountsTopicV1,
			Key:   []byte(account.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "event_id", Value: []byte(eventID.String())}, // Outbox will fill this
				{Key: "event_type", Value: []byte(AccountCreatedEvent)},
				{Key: "event_version", Value: []byte("1")},
				{Key: "producer", Value: []byte(contracts.SsoSvcProducer)},
				{Key: "content_type", Value: []byte("application/json")},
			},
		},
	)

	return err
}
