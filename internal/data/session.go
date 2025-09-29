package data

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

type Sessions interface {
	Insert(ctx context.Context, input models.SessionRow) error

	Update(ctx context.Context) error
	UpdateLastUsed(lastUsed time.Time) Sessions
	UpdateToken(token string) Sessions

	Delete(ctx context.Context) error
	Select(ctx context.Context) ([]models.SessionRow, error)
	Get(ctx context.Context) (models.SessionRow, error)

	FilterID(id uuid.UUID) Sessions
	FilterUserID(userID uuid.UUID) Sessions

	Page(limit, offset uint) Sessions
	Count(ctx context.Context) (uint, error)

	OrderCreatedAt(ascending bool) Sessions

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
