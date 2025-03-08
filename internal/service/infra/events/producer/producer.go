package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/evebody"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/kafig"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer interface {
	AccountCreate(account models.Account) error
}

type events struct {
	cfg config.Config
	log *logrus.Logger
}

func NewBroker(cfg *config.Config, log *logrus.Logger) Producer {
	return &events{
		cfg: *cfg,
		log: log,
	}
}

type InternalEvent struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

func (e *events) AccountCreate(account models.Account) error {
	writer := kafka.Writer{
		Addr:         kafka.TCP(e.cfg.Kafka.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 0,
		Async:        false,
		RequiredAcks: kafka.RequireAll,
	}

	body, err := json.Marshal(evebody.AccountCreated{
		AccountID: account.ID.String(),
		Email:     account.Email,
		Role:      string(account.Role),
		Timestamp: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal AccountCreated event: %w", err)
	}

	evt := InternalEvent{
		EventType: "account_created",
		Data:      body,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal AccountCreated event: %w", err)
	}

	msg := kafka.Message{
		Topic: kafig.AccountsTopic,
		Key:   []byte(account.ID.String()),
		Value: data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

type TopicConfig struct {
	Topic    string
	Callback func(ctx context.Context, m kafka.Message, evt InternalEvent) error
}
