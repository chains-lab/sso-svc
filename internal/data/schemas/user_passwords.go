package schemas

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UsersPasswordQ interface {
	Insert(ctx context.Context, input UserPasswordModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]UserPasswordModel, error)
	Get(ctx context.Context) (UserPasswordModel, error)

	FilterID(id uuid.UUID) UsersPasswordQ

	Update(ctx context.Context, input UserPassUpdateInput) error

	Page(limit, offset uint) UsersPasswordQ
	Count(ctx context.Context) (uint, error)
}

type UserPasswordModel struct {
	ID        uuid.UUID `db:"user_id"`
	PassHash  string    `db:"password_hash"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserPassUpdateInput struct {
	PassHash  *string
	UpdatedAt time.Time
}
