package user

import (
	"context"

	userProto "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/google/uuid"
)

type App interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	Login(ctx context.Context, email string, client string) (models.Session, models.TokensPair, error)

	AdminCreateUser(ctx context.Context, initiatorID uuid.UUID, email string, input app.AdminCreateUserInput) (models.User, error)

	UpdateUserVerified(ctx context.Context, initiatorID, userID uuid.UUID, verified bool) (models.User, error)
	UpdateUserSuspended(ctx context.Context, initiatorID, userID uuid.UUID, suspended bool) (models.User, error)
}

type Service struct {
	app App
	cfg config.Config

	userProto.UnimplementedUserServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
