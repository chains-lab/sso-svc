package producer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	addr   string
	outbox outbox
	log    logium.Logger
}

type outbox interface {
	CreateOutboxEvent(
		ctx context.Context,
		event contracts.Message,
	) error

	GetPendingOutboxEvents(
		ctx context.Context,
		limit int32,
	) ([]contracts.OutboxEvent, error)

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

func New(log logium.Logger, addr string, outbox outbox) *Service {
	return &Service{
		addr:   addr,
		outbox: outbox,
		log:    log,
	}
}

func (s Service) Publish(
	ctx context.Context,
	event contracts.Message,
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

const eventOutboxRetryDelay = 1 * time.Minute

func (s Service) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := s.outbox.GetPendingOutboxEvents(ctx, 100)
			if err != nil {

				continue
			}

			var sentIDs []uuid.UUID
			var NonSentIDs []uuid.UUID
			for _, eventData := range events {
				err = s.Publish(ctx, eventData.ToEventData())
				if err != nil {
					NonSentIDs = append(NonSentIDs, eventData.ID)
					s.log.Printf("outbox: publish event ID %s: %v", eventData.ID, err)
					continue
				}
				sentIDs = append(sentIDs, eventData.ID)
			}

			if len(sentIDs) > 0 {
				err = s.outbox.MarkOutboxEventsSent(ctx, sentIDs)
				if err != nil {
					s.log.Printf("outbox: mark events as sent: %v", err)
				}
			}

			if len(NonSentIDs) > 0 {
				err = s.outbox.DelayOutboxEvents(ctx, NonSentIDs, eventOutboxRetryDelay)
				if err != nil {
					s.log.Printf("outbox: delay events: %v", err)
				}
			}
		}
	}
}
