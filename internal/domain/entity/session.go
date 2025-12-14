package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Session) IsNil() bool {
	return s.ID == uuid.Nil
}

type SessionsCollection struct {
	Data  []Session `json:"repo"`
	Page  int32     `json:"page"`
	Size  int32     `json:"size"`
	Total int64     `json:"total"`
}
