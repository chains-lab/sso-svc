package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionsCollection struct {
	Data  []Session `json:"data"`
	Page  uint      `json:"page"`
	Size  uint      `json:"size"`
	Total uint      `json:"total"`
}
