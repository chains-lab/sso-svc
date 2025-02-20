package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/tokens/identity"
)

type Account struct {
	ID        uuid.UUID
	Email     string
	Role      identity.IdnType
	UpdatedAt time.Time
	CreatedAt time.Time
}
