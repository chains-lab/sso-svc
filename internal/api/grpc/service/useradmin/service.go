package useradmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	adminPoroto "github.com/chains-lab/sso-proto/gen/go/svc/useradmin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type App interface {
	AdminCreateUser(ctx context.Context, initiatorID uuid.UUID, email string, input app.AdminCreateUserInput) (models.User, error)

	UpdateUserVerified(ctx context.Context, userID uuid.UUID, verified bool) (models.User, error)
	UpdateUserSuspended(ctx context.Context, userID uuid.UUID, suspended bool) (models.User, error)

	GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error)

	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
}

type Service struct {
	app App
	cfg config.Config

	adminPoroto.UnimplementedUserAdminServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func (s Service) ComparisonRightsForAdmins(ctx context.Context, initiatorStrID, userStrID string) (initiator, user models.User, err error) {
	initiatorID, err := uuid.Parse(initiatorStrID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return initiator, user, problems.UnauthenticatedError(ctx, fmt.Sprintf("invalid initiator ID format: %s", initiatorStrID))
	}

	userID, err := uuid.Parse(userStrID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return initiator, user, problems.InvalidArgumentError(
			ctx,
			fmt.Sprintf("invalid user ID format: %s", userStrID),
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid format user id",
			})
	}

	initiator, err = s.app.GetUserByID(ctx, initiatorID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorUserNotFound):
			return initiator, user, problems.UnauthenticatedError(ctx, fmt.Sprintf("initiator user %s not found", initiatorStrID))
		}

		return initiator, user, err
	}

	if initiator.Suspended {
		return initiator, user, errx.RaiseInitiatorUserSuspended(
			ctx,
			fmt.Errorf("initiator %s is suspended", initiatorStrID),
			initiatorStrID,
		)
	}

	user, err = s.app.GetUserByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return initiator, user, errx.RaiseInitiatorRoleIsLowThanTarget(
				ctx,
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}
