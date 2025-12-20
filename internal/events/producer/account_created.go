package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/kafkakit/header"
	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/internal/events/contracts"
)

func (s Service) WriteAccountCreated(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	payload, err := json.Marshal(contracts.AccountCreatedPayload{
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
				{Key: header.EventID, Value: []byte(eventID.String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.AccountUsernameChangeEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.SsoSvcProducer)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
