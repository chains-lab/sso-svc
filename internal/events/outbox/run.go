package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
