package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	LastUsed  time.Time `db:"last_used"`
	CreatedAt time.Time `db:"created_at"`
}
