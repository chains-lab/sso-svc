package session

import (
	"context"

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

	GetSession(ctx context.Context, sessionID uuid.UUID) (schemas.Session, error)
	GetOneSessionForUser(ctx context.Context, userID, sessionID uuid.UUID) (schemas.Session, error)
	GetAllSessionsForUser(ctx context.Context, userID uuid.UUID, limit, offset uint) ([]schemas.Session, uint, error)

	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteOneSessionForUser(ctx context.Context, userID, sessionID uuid.UUID) error
	DeleteAllSessionsForUser(ctx context.Context, userID uuid.UUID) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (schemas.User, error)
	GetUserByEmail(ctx context.Context, email string) (schemas.User, error)
}
