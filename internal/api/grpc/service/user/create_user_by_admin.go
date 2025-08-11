package user

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateUserByAdmin(ctx context.Context, req *svc.CreateUserByAdminRequest) (*svc.User, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access: only admin or super admin can create user")

		return nil, problem.PermissionDeniedError(ctx, "only admins roles create user")
	}

	userRole, err := roles.ParseRole(req.Role)
	if err != nil {
		return nil, problem.InvalidArgumentError(ctx, "user role is not allowed", &errdetails.BadRequest_FieldViolation{
			Field:       "role",
			Description: "invalid role",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	user, err := s.app.AdminCreateUser(ctx, initiatorID, req.Role, app.AdminCreateUserInput{
		Role:     userRole,
		Verified: req.Verified,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create admin user")

		return nil, err
	}

	logger.Log(ctx).Warnf("admin user %s created successfully", user.ID)
	return response.User(user), nil
}
