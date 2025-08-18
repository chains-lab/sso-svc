package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	EmailVer       bool      `json:"email_verified"`
	EmailUpdatedAt time.Time `json:"email_updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}
