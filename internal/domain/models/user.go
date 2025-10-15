package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID `json:"id"`
	SysRole string    `json:"role"`
	Status  string    `json:"status"`

	CityID      *uuid.UUID `json:"city_id,omitempty"`
	CityRole    *string    `json:"city_role,omitempty"`
	CompanyID   *uuid.UUID `json:"company_id,omitempty"`
	CompanyRole *string    `json:"company_role,omitempty"`

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
