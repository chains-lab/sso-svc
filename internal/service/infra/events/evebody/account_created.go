package evebody

import (
	"time"
)

type AccountCreated struct {
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
