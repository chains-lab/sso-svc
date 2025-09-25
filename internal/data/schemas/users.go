package schemas

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UsersQ interface {
	Insert(ctx context.Context, input UserModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]UserModel, error)
	Get(ctx context.Context) (UserModel, error)

	FilterID(id uuid.UUID) UsersQ
	FilterRole(role string) UsersQ

	Update(ctx context.Context, input UserUpdateInput) error

	Page(limit, offset uint) UsersQ
	Count(ctx context.Context) (uint, error)
}

type UserModel struct {
	ID        uuid.UUID `db:"id"`
	Role      string    `db:"role"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

type UserUpdateInput struct {
	Status *string
}
