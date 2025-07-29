package models

import (
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Role      roles.Role `json:"role"`
	Verified  bool       `json:"verified,omitempty"`
	Suspended bool       `json:"suspended,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}
