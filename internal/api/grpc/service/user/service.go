package user

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

type Service struct {
	app *app.App
	cfg config.Config

	userProto.UnimplementedAuthServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func (s Service) allowedAccessForUser(ctx context.Context, initiatorID, userID uuid.UUID) (bool, error) {
	initiator, err := s.app.GetInitiatorByID(ctx, initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error fetching user with ID %s", initiatorID)

		return false, err
	}

	if initiator.Role == roles.SuperUser {
		return true, nil
	}

	user, err := s.app.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error fetching user with ID %s", userID)

		return false, err
	}

	if roles.CompareRolesUser(initiator.Role, user.Role) > 0 {
		return true, nil
	}

	return false, nil
}
