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

type Kafka interface {
	AccountCreate(account models.Account) error

	sendMessage(msg kafka.Message) error
}

type broker struct {
	Writer *kafka.Writer
	cfg    *config.Config
	log    *logrus.Logger
}

type TopicConfig struct {
	Topic    string
	Callback func(ctx context.Context, message kafka.Message) error
}

func NewBroker(cfg *config.Config, log *logrus.Logger) Kafka {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 0,
		Async:        false,
		RequiredAcks: kafka.RequireAll,
	}

	return &broker{
		Writer: writer,
		cfg:    cfg,
		log:    log,
	}
}

func (b *broker) sendMessage(msg kafka.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.Writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

type internalEvent struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data"`
}

func (b *broker) AccountCreate(account models.Account) error {
	body := evebody.AccountCreated{
		AccountID: account.ID.String(),
		Email:     account.Email,
		Role:      string(account.Role),
		Timestamp: time.Now(),
	}

	evt := internalEvent{
		EventType: "account_created",
		Data:      body,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal AccountCreated event: %w", err)
	}

	msg := kafka.Message{
		Topic: kafig.AccountsTopic,
		Key:   []byte(body.AccountID),
		Value: data,
	}

	return b.sendMessage(msg)
}
