package publisher

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	addr string
}

func New(addr string) *Service {
	return &Service{
		addr: addr,
	}
}

type Envelope interface {
	MarshalJSON() ([]byte, error)
	EventType() string
	EventVersion() string
	EventTime() time.Time
}

func (p Service) publish(
	ctx context.Context,
	topic, key string,
	envelope Envelope,
	headers ...kafka.Header,
) error {
	body, err := envelope.MarshalJSON()
	if err != nil {
		return err
	}

	writer := kafka.Writer{
		Addr:         kafka.TCP(p.addr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Compression:  kafka.Snappy,
		BatchTimeout: 50 * time.Millisecond,
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("kafka: close publisher: %v", err)
		}
	}()

	msg := kafka.Message{
		Key:   []byte(key),
		Value: body,
		Time:  envelope.EventTime(),
		Headers: append(headers,
			kafka.Header{Key: "event_type", Value: []byte(envelope.EventType())},
			kafka.Header{Key: "event_version", Value: []byte(envelope.EventVersion())},
			kafka.Header{Key: "content_type", Value: []byte("application/json")},
		),
	}

	return writer.WriteMessages(ctx, msg)
}
