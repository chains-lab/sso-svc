package events

import "time"

type RoleUpdated struct {
	Event     string    `json:"event"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
