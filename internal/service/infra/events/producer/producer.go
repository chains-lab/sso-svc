package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	AccountCreate(account models.Account) error
}

type producer struct {
	brokers net.Addr
	writer  *kafka.Writer
}

func NewProducer(cfg *config.Config) Producer {
	return &producer{
		brokers: kafka.TCP(cfg.Kafka.Brokers...),
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    1,
			BatchTimeout: 0,
			Async:        false,
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *producer) AccountCreate(account models.Account) error {
	body, err := json.Marshal(events.AccountCreated{
		AccountID: account.ID.String(),
		Email:     account.Email,
		Role:      string(account.Role),
		Timestamp: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal AccountCreated event: %w", err)
	}

	return p.sendMessage(events.AccountCreateTopic, events.SubscriptionDeactivatedType, account.ID.String(), body)
}

func (p *producer) sendMessage(topic string, event string, key string, body []byte) error {
	evt := events.InternalEvent{
		EventType: event,
		Data:      body,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription activate event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Value: data,
		Key:   []byte(key),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
