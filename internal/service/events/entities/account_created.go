package entities

import (
	"time"
)

type AccountCreated struct {
	Event     string    `json:"event"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
