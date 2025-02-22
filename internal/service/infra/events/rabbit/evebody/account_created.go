package evebody

import (
	"time"
)

type AccountCreated struct {
	Event     string    `json:"event"`
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
