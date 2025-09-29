package domain

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/google/uuid"
)

type SessionSvc interface {
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error

	Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	GetForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)

	ListForUser(
		ctx context.Context,
		userID uuid.UUID,
		page uint,
		size uint,
	) (models.SessionsCollection, error)

	Login(ctx context.Context, email, password string) (models.TokensPair, error)
	LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error)
	Create(ctx context.Context, userID uuid.UUID, role string) (models.TokensPair, error)

	Refresh(ctx context.Context, oldRefreshToken string) (models.TokensPair, error)
}
