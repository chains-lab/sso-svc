package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID           uuid.UUID
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      json.RawMessage
	Status       string
	Attempts     int32
	NextRetryAt  time.Time
	CreatedAt    time.Time
	SentAt       *time.Time
}
