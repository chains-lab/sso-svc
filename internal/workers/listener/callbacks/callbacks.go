package callbacks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events"
	"github.com/segmentio/kafka-go"
)

func SubscriptionStatus(ctx context.Context, svc *service.Service, m kafka.Message, evt events.InternalEvent) error {
	var ps events.SubscriptionActivated
	if err := json.Unmarshal(evt.Data, &ps); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}
	userID, err := uuid.Parse(string(m.Key))
	if err != nil {
		return fmt.Errorf("failed to parse payment evemt key: %w", err)
	}
	typeID, err := uuid.Parse(ps.TypeID)
	if err != nil {
		return fmt.Errorf("failed to parse payment evemt key: %w", err)
	}
	if evt.EventType == events.SubscriptionActivatedType {
		return svc.Domain.SubscriptionUpdate(ctx, userID, &typeID)
	}
	return svc.Domain.SubscriptionUpdate(ctx, userID, nil)
}
