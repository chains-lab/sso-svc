package contracts

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      interface{}
}

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

func (e OutboxEvent) ToEventData() Message {
	return Message{
		Topic:        e.Topic,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		Key:          e.Key,
		Payload:      e.Payload,
	}
}
