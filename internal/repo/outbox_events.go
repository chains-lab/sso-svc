package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/chains-lab/sso-svc/internal/events/outbox"
	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

type CreateOutboxEventParams struct {
	ID           uuid.UUID
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      json.RawMessage
}

func (r *Repository) CreateOutboxEvent(
	ctx context.Context,
	params CreateOutboxEventParams,
) (outbox.OutboxEvent, error) {
	res, err := r.sql.CreateOutboxEvent(ctx, pgdb.CreateOutboxEventParams{
		ID:           params.ID,
		Topic:        params.Topic,
		EventType:    params.EventType,
		EventVersion: params.EventVersion,
		Key:          params.Key,
		Payload:      params.Payload,
	})
	if err != nil {
		return outbox.OutboxEvent{}, err
	}

	return res.ToModel(), nil
}

func (r *Repository) GetPendingOutboxEvents(
	ctx context.Context,
	limit int32,
) ([]outbox.OutboxEvent, error) {
	res, err := r.sql.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		return nil, err
	}

	events := make([]outbox.OutboxEvent, len(res))
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
