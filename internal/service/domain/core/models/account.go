package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
)

type Account struct {
	ID        uuid.UUID
	Email     string
	Role      roles.UserRole
	UpdatedAt time.Time
	CreatedAt time.Time
}
