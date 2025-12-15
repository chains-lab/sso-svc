package outbox

import (
	"context"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	writer writer
	repo   repository
	log    logium.Logger
}

func New(log logium.Logger, writer writer, repo repository) Service {
	return Service{
		repo:   repo,
		writer: writer,
		log:    log,
	}
}

type writer interface {
	Write(
		ctx context.Context,
		event contracts.Event,
		headers ...kafka.Header,
	) error
}

type repository interface {
	GetPendingOutboxEvents(
		ctx context.Context,
		limit int32,
	) ([]OutboxEvent, error)

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

const eventRetryDelay = 1 * time.Minute

func (o Service) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := o.repo.GetPendingOutboxEvents(ctx, 100)
			if err != nil {

				continue
			}

			var sentIDs []uuid.UUID
			var NonSentIDs []uuid.UUID
			for _, eventData := range events {
				err = o.writer.Write(ctx, eventData.ToEventData())
				if err != nil {
					NonSentIDs = append(NonSentIDs, eventData.ID)
					o.log.Printf("outbox: publish event ID %s: %v", eventData.ID, err)
					continue
				}
				sentIDs = append(sentIDs, eventData.ID)
			}

			if len(sentIDs) > 0 {
				err = o.repo.MarkOutboxEventsSent(ctx, sentIDs)
				if err != nil {
					o.log.Printf("outbox: mark events as sent: %v", err)
				}
			}

			if len(NonSentIDs) > 0 {
				err = o.repo.DelayOutboxEvents(ctx, NonSentIDs, eventRetryDelay)
				if err != nil {
					o.log.Printf("outbox: delay events: %v", err)
				}
			}
		}
	}
}
