package events

import (
	"time"
)

type AccountCreated struct {
	Event     string    `json:"event"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}
