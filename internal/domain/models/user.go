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
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) IsNil() bool {
	return u.ID == uuid.Nil
}

type UserPassword struct {
	Hash      string    `json:"hash"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserPassword) IsNil() bool {
	return u.Hash == ""
}
