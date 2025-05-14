package models

import (
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         roles.Role `json:"role"`
	Subscription uuid.UUID  `json:"subscription"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
