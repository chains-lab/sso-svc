package eventlistener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/events"
	"github.com/chains-lab/chains-auth/internal/events/reader"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Subscription interface {
	SubscriptionUpdate(ctx context.Context, AccountID uuid.UUID, subscriptionID uuid.UUID) error
}

type Listener struct {
	subscription Subscription

	cfg *config.Config
	log *logrus.Entry
}

func NewListener(cfg *config.Config, app *app.App, logger *logrus.Logger) *Listener {
	return &Listener{
		subscription: app,
		cfg:          cfg,
		log:          logger.WithField("module", "event-listener"),
	}
}

func (l *Listener) Listen(ctx context.Context, cfg *config.Config) {
	subscriptionReader := reader.NewReader(l.log, kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          events.SubscriptionsTopic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	}))

	subscriptionChanel := subscriptionReader.ListenChan(ctx)

	go func(ctx context.Context) {
		for event := range subscriptionChanel {
			var eve events.AccountCreated
			if err := json.Unmarshal(event.Data, &eve); err != nil {
				l.log.WithError(err).Error("Error unmarshalling account create event")
				continue
			}

			switch event.EventType {
			case events.SubscriptionActivateType:
				var subEvent events.SubscriptionEvent
				if err := json.Unmarshal(event.Data, &subEvent); err != nil {
					l.log.WithError(err).Error("failed to unmarshal SubscriptionActivateType")
					continue
				}

				err := l.subscription.SubscriptionUpdate(ctx, subEvent.AccountID, subEvent.SubscriptionID)
				if err != nil {
					l.log.WithError(err).Error("error processing subscription activate event")
					continue
				}
			case events.SubscriptionDeactivateType:
				var subEvent events.SubscriptionEvent
				if err := json.Unmarshal(event.Data, &subEvent); err != nil {
					l.log.WithError(err).Error("failed to unmarshal SubscriptionActivateType")
					continue
				}

				err := l.subscription.SubscriptionUpdate(ctx, subEvent.AccountID, subEvent.SubscriptionID)
				if err != nil {
					l.log.WithError(err).Error("error processing subscription activate event")
					continue
				}
			default:
				l.log.WithField("event", event).Error("Unknown event type")
			}
		}
	}(ctx)

	<-ctx.Done()
	l.log.Info("Producer listener stopped")
}
