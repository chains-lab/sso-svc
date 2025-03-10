package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/tokens/identity"
)

type Account struct {
	ID           uuid.UUID        `json:"id"`
	Email        string           `json:"email"`
	Role         identity.IdnType `json:"role"`
	Subscription *uuid.UUID       `json:"subscription,omitempty"`
	UpdatedAt    time.Time        `json:"updated_at"`
	CreatedAt    time.Time        `json:"created_at"`
}
