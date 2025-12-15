package writer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/segmentio/kafka-go"
)

func (s Service) Write(
	ctx context.Context,
	event contracts.Event,
	headers ...kafka.Header,
) error {
	writer := kafka.Writer{
		Addr:         kafka.TCP(s.addr),
		Topic:        event.Topic,
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

	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.Key),
		Value: payload,
		Time:  time.Now().UTC(),
		Headers: append(headers,
			kafka.Header{Key: "event_type", Value: []byte(event.EventType)},
			kafka.Header{Key: "event_version", Value: []byte(strconv.Itoa(int(event.EventVersion)))},
			kafka.Header{Key: "content_type", Value: []byte("application/json")},
		),
	}

	return writer.WriteMessages(ctx, msg)
}
