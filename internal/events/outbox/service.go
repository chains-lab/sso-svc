package outbox

import (
	"context"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const eventRetryDelay = 1 * time.Minute

type Service struct {
	writer writer
	repo   repository
	log    logium.Logger
}

func New(log logium.Logger, writer writer, repo repository) Service {
	return Service{
		repo:   repo,
		writer: writer,
		log:    log,
	}
}

type writer interface {
	Write(
		ctx context.Context,
		event contracts.Event,
		headers ...kafka.Header,
	) error
}

type repository interface {
	GetPendingOutboxEvents(
		ctx context.Context,
		limit int32,
	) ([]EventData, error)

	MarkOutboxEventsSent(
		ctx context.Context,
		ids []uuid.UUID,
	) error

	DelayOutboxEvents(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) error
}

type EventData struct {
	ID           uuid.UUID
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      interface{}
	Status       string
	Attempts     int32
	NextRetryAt  time.Time
	CreatedAt    time.Time
	SentAt       *time.Time
}

func (o EventData) ToEventData() contracts.Event {
	return contracts.Event{
		Topic:        o.Topic,
		EventType:    o.EventType,
		EventVersion: o.EventVersion,
		Key:          o.Key,
		Payload:      o.Payload,
	}
}
