package events

import (
	"encoding/json"
	"time"
)

type Envelope[T any] struct {
	Event     string    `json:"event"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}

func (e Envelope[T]) MarshalJSON() ([]byte, error) {
	type alias Envelope[T]
	return json.Marshal(alias(e))
}

func (e Envelope[T]) EventType() string {
	return e.Event
}

func (e Envelope[T]) EventVersion() string {
	return e.Version
}

func (e Envelope[T]) EventTime() time.Time {
	return e.Timestamp
}

func (e Envelope[T]) EventData() interface{} {
	return e.Data
}
