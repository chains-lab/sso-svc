package events

import (
	"encoding/json"
	"time"
)

const (
	AccountsTopic      = "accounts"
	AccountCreateTopic = "account_create"
	SubscriptionsTopic = "subscriptions"

	// events types

	AccountCreateEventType      = "account_create"
	SubscriptionActivatedType   = "subscription_activated"
	SubscriptionDeactivatedType = "subscription_deactivated"
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

type SubscriptionActivated struct {
	AccountID string    `json:"account_id"`
	TypeID    string    `json:"type_id"`
	PlanID    string    `json:"plan_id"`
	Timestamp time.Time `json:"timestamp"`
}

type SubscriptionDeactivated struct {
	AccountID string    `json:"account_id"`
	TypeID    string    `json:"type_id"`
	PlanID    string    `json:"plan_id"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountRoleUpdated struct {
	Event     string    `json:"event"`
	AccountID string    `json:"account_id"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
