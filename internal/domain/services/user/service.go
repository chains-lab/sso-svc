package user

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

func New(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (schemas.User, error)
	GetUserByEmail(ctx context.Context, email string) (schemas.User, error)

	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string, updatedAt time.Time) error
	UpdateUserEmailVerification(ctx context.Context, userID uuid.UUID, verified bool, updatedAt time.Time) error

	DeleteUser(ctx context.Context, userID uuid.UUID) error

	DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error
}
