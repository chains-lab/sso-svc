package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/chains-lab/sso-svc/internal/events/outbox"
	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

type CreateOutboxEventParams struct {
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      json.RawMessage
}

func (r *Repository) CreateOutboxEvent(
	ctx context.Context,
	event contracts.Event,
) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	_, err = r.sql.CreateOutboxEvent(ctx, pgdb.CreateOutboxEventParams{
		Topic:        event.Topic,
		EventType:    event.EventType,
		EventVersion: event.EventVersion,
		Key:          event.Key,
		Payload:      payloadBytes,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPendingOutboxEvents(
	ctx context.Context,
	limit int32,
) ([]outbox.EventData, error) {
	res, err := r.sql.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []outbox.EventData{}, nil
		}
		return nil, err
	}

	events := make([]outbox.EventData, len(res))
	for i, e := range res {
		events[i] = e.ToModel()
	}

	return events, nil
}

func (r *Repository) MarkOutboxEventsSent(
	ctx context.Context,
	ids []uuid.UUID,
) error {
	res := r.sql.MarkOutboxEventsSent(ctx, pgdb.MarkOutboxEventsSentParams{
		Column1: ids,
		SentAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
	})

	return res
}

func (r *Repository) DelayOutboxEvents(
	ctx context.Context,
	ids []uuid.UUID,
	delay time.Duration,
) error {
	res := r.sql.DelayOutboxEvents(ctx, pgdb.DelayOutboxEventsParams{
		Column1:     ids,
		NextRetryAt: time.Now().Add(delay),
	})

	return res
}
