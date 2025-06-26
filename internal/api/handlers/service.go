package handlers

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/interceptors"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/utils/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type App interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	GetSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)
	SelectUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	Login(ctx context.Context, email string, role roles.Role, client string) (models.Session, models.TokensPair, error)
	Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, error)

	DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error

	TerminateUserSessions(ctx context.Context, userID uuid.UUID) error

	AdminUpdateUserRole(ctx context.Context, initiatorID, userID uuid.UUID, role roles.Role) (models.User, error)
	AdminUpdateUserSubscription(ctx context.Context, initiatorID, userID, subscriptionID uuid.UUID) (models.User, error)
	AdminUpdateUserVerified(ctx context.Context, initiatorID, userID uuid.UUID, verified bool) (models.User, error)
	AdminUpdateUserSuspended(ctx context.Context, initiatorID, userID uuid.UUID, suspended bool) (models.User, error)

	AdminTerminateSessions(ctx context.Context, initiatorID, userID uuid.UUID) error
	AdminDeleteUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error
}

type Service struct {
	app App
	cfg config.Config

	svc.ServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func Log(ctx context.Context, requestID uuid.UUID) *logrus.Entry {
	entry, ok := ctx.Value(interceptors.LogCtxKey).(*logrus.Entry)
	if !ok {
		entry = logrus.NewEntry(logrus.New())
	}
	return entry.WithField("request_id", requestID)
}

func Meta(ctx context.Context) interceptors.MetaData {
	md, ok := ctx.Value(interceptors.MetaCtxKey).(interceptors.MetaData)
	if !ok {
		return interceptors.MetaData{}
	}
	return md
}
