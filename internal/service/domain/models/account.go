package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID
	Email     string
	Role      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
