package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chains-lab/chains-auth/internal/events/bodies"
	"github.com/chains-lab/chains-auth/internal/utils/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type CreateUserWriters struct {
	writer kafka.Writer
}

func NewUserCreateWriters(cfg config.Config, logger *logrus.Entry) *CreateUserWriters {
	w := &CreateUserWriters{
		writer: kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Topic:        bodies.UserCreateTopic,
			RequiredAcks: kafka.RequireOne,
			Async:        false,
			Logger:       kafka.LoggerFunc(logger.Infof),
			ErrorLogger:  kafka.LoggerFunc(logger.Errorf),
		},
	}
	return w
}

func (w *CreateUserWriters) CreateUser(ctx context.Context, created bodies.UserCreated) error {
	dataPayload, err := json.Marshal(created)
	if err != nil {
		return fmt.Errorf("marshal UserCreated: %w", err)
	}

	evt := bodies.InternalEvent{
		Type: bodies.UserCreateType,
		Data: json.RawMessage(dataPayload),
	}

	msgValue, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal InternalEvent: %w", err)
	}

	key := created.UserID

	msg := kafka.Message{
		Key:   []byte(key),
		Value: msgValue,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := w.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write kafka message: %w", err)
	}
	return nil
}
