package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/tokens/identity"
)

type Account struct {
	ID           uuid.UUID     `json:"id"`
	Email        string        `json:"email"`
	Role         identity.Role `json:"role"`
	Subscription uuid.UUID     `json:"subscription"`
	UpdatedAt    time.Time     `json:"updated_at"`
	CreatedAt    time.Time     `json:"created_at"`
}
