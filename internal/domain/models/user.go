package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	Email     string    `json:"email"`
	EmailVer  bool      `json:"email_verified"`
	CreatedAt time.Time `json:"created_at"`
}
