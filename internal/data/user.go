package data

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

type Users interface {
	Insert(ctx context.Context, input models.UserRow) error
	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]models.UserRow, error)
	Get(ctx context.Context) (models.UserRow, error)

	FilterID(id uuid.UUID) Users
	FilterRole(role string) Users
	FilterEmail(email string) Users

	Update(ctx context.Context, updateAt time.Time) error
	UpdateEmail(email string) Users
	UpdateStatus(status string) Users
	UpdatePassword(passwordHash string, passwordUpAt time.Time) Users
	UpdateEmailVerified(emailVer bool) Users

	Page(limit, offset uint) Users
	Count(ctx context.Context) (uint, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
