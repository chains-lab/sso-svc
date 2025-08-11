package user

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateUserSuspended(ctx context.Context, req *svc.UpdateUserSuspendedRequest) (*svc.User, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access: only admin or super admin can update user suspended status")

		return nil, problem.PermissionDeniedError(
			ctx,
			"only admin or super admin can update user suspended status",
		)
	}

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse user ID")

		return nil, problem.InvalidArgumentError(
			ctx,
			"invalid user ID format",
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "user ID must be a valid UUID",
			})
	}

	user, err := s.app.UpdateUserSuspended(ctx, initiatorID, userID, req.Suspended)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user suspended status")

		return nil, err
	}

	logger.Log(ctx).Warnf("user %s suspended status updated to %v successfully", user.ID, req.Suspended)

	return response.User(user), nil
}
