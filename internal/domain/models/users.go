package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	EmailVer  bool      `json:"email_verified"`
	CreatedAt time.Time `json:"created_at"`
}
