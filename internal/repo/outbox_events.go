package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
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
	event contracts.Message,
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
) ([]contracts.OutboxEvent, error) {
	res, err := r.sql.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []contracts.OutboxEvent{}, nil
		}
		return nil, err
	}

	events := make([]contracts.OutboxEvent, len(res))
	for i, e := range res {
		events[i] = e.ToEntity()
	}

	return events, nil
}

func (r *Repository) MarkOutboxEventsSent(
	ctx context.Context,
	ids []uuid.UUID,
) error {
	return r.sql.MarkOutboxEventsSent(ctx, pgdb.MarkOutboxEventsSentParams{
		Column1: ids,
		SentAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
	})
}

func (r *Repository) DelayOutboxEvents(
	ctx context.Context,
	ids []uuid.UUID,
	delay time.Duration,
) error {
	return r.sql.DelayOutboxEvents(ctx, pgdb.DelayOutboxEventsParams{
		Column1:     ids,
		NextRetryAt: time.Now().Add(delay),
	})
}
