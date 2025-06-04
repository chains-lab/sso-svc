package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`

	Access  string `json:"token"`
	Refresh string `json:"refresh_at"`

	Client    string    `json:"client"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}
