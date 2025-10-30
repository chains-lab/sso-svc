package subscriber

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	reader *kafka.Reader
}

func New(addr string, topic string, groupID string) *Service {
	cfg := kafka.ReaderConfig{
		Brokers:         strings.Split(addr, ","),
		Topic:           topic,
		GroupID:         groupID,
		MinBytes:        1e3,  // 1KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         500 * time.Millisecond,
		ReadLagInterval: -1,
	}

	if groupID == "" {
		cfg.StartOffset = kafka.FirstOffset
	}

	return &Service{reader: kafka.NewReader(cfg)}
}

func (s *Service) Subscribe(
	ctx context.Context,
	filterEventType string,
	callback func(ctx context.Context, event kafka.Message) error,
) error {
	defer func() {
		if err := s.reader.Close(); err != nil {
			log.Printf("failed to close reader: %v", err)
		}
	}()

	for {
		m, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("read error: %v", err)
			continue
		}

		if filterEventType != "" {
			ok := false
			for _, h := range m.Headers {
				if h.Key == "event_type" && string(h.Value) == filterEventType {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("callback panic recovered: %v", r)
				}
			}()
			if err := callback(ctx, m); err != nil {
				log.Printf("callback error (%s): %v", filterEventType, err)
			}
		}()
	}
}
