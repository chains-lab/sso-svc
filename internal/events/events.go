package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	AccountsTopic      = "accounts"
	SubscriptionsTopic = "subscriptions"

	AccountCreateType = "ACCOUNT_CREATE"

	SubscriptionActivateType   = "SUBSCRIPTION_ACTIVATE"
	SubscriptionDeactivateType = "SUBSCRIPTION_DEACTIVATE"
)

type InternalEvent struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type AccountCreated struct {
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type SubscriptionEvent struct {
	AccountID      uuid.UUID `json:"account_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Timestamp      time.Time `json:"timestamp"`
}
