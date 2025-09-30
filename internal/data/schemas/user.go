package schemas

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID     uuid.UUID `db:"id"`
	Role   string    `db:"role"`
	Status string    `db:"status"`

	PasswordHash string    `db:"password_hash"`
	PasswordUpAt time.Time `db:"password_updated_at"`

	Email    string `db:"email"`
	EmailVer bool   `db:"email_verified"`

	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}
