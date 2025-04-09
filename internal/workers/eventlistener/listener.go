package eventlistener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/events"
	"github.com/hs-zavet/sso-oauth/internal/events/reader"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	cfg *config.Config
	app *app.App
	log *logrus.Entry
}

func NewListener(cfg *config.Config, app *app.App, logger *logrus.Logger) *Listener {
	return &Listener{
		cfg: cfg,
		app: app,
		log: logger.WithField("module", "event-listener"),
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
				//l.app.SubscriptionActivate(ctx, event)
			case events.SubscriptionDeactivateType:
				//l.app.AccountDeactivate(ctx, event)
			default:
				l.log.WithField("event", event).Error("Unknown event type")
			}
		}
	}(ctx)

	<-ctx.Done()
	l.log.Info("Producer listener stopped")
}
