package schemas

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionsQ interface {
	Insert(ctx context.Context, input Session) error
	Update(ctx context.Context, input UpdateSessionInput) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]Session, error)
	Get(ctx context.Context) (Session, error)

	FilterID(id uuid.UUID) SessionsQ
	FilterUserID(userID uuid.UUID) SessionsQ

	Page(limit, offset uint) SessionsQ
	Count(ctx context.Context) (uint, error)

	OrderCreatedAt(ascending bool) SessionsQ
}

type Session struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	LastUsed  time.Time `db:"last_used"`
	CreatedAt time.Time `db:"created_at"`
}

type UpdateSessionInput struct {
	Token    *string
	LastUsed time.Time
}
