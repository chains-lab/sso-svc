package outbox

import (
	"time"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type OutboxEvent struct {
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

func (o OutboxEvent) ToEventData() contracts.Event {
	return contracts.Event{
		Topic:        o.Topic,
		EventType:    o.EventType,
		EventVersion: o.EventVersion,
		Key:          o.Key,
		Payload:      o.Payload,
	}
}
