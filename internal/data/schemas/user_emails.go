package schemas

import (
	"context"

	"github.com/google/uuid"
)

type UsersEmailQ interface {
	Insert(ctx context.Context, input UserEmailModel) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]UserEmailModel, error)
	Get(ctx context.Context) (UserEmailModel, error)

	FilterID(id uuid.UUID) UsersEmailQ
	FilterEmail(email string) UsersEmailQ

	Update(ctx context.Context, input UserEmailUpdateInput) error

	Page(limit, offset uint) UsersEmailQ
	Count(ctx context.Context) (uint, error)
}

type UserEmailModel struct {
	ID       uuid.UUID `db:"user_id"`
	Email    string    `db:"email"`
	Verified bool      `db:"verified"`
}

type UserEmailUpdateInput struct {
	Email    *string
	Verified *bool
}
