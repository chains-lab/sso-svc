package evebody

import "time"

type AccountRoleUpdated struct {
	Event     string    `json:"event"`
	AccountID string    `json:"account_id"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
