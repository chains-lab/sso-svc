package session

import (
	"context"

	sessPoroto "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
)

type App interface {
	GetSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)
	GetSessions(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error)

	Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, error)

	DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error
	DeleteSessions(ctx context.Context, userID uuid.UUID) error

	AdminCreateUser(ctx context.Context, input app.AdminCreateUserInput) (models.User, error)

	AdminDeleteSessions(ctx context.Context, initiatorID, userID uuid.UUID) error
	AdminDeleteSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error
}

type Service struct {
	app App
	cfg config.Config

	sessPoroto.UnimplementedSessionServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
