package session

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteSessionsByAdmin(ctx context.Context, req *svc.DeleteSessionsByAdminRequest) (*emptypb.Empty, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access: only admin or super admin can create user")

		return nil, problems.PermissionDeniedError(ctx, "only admins or super admin can create user")
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, problems.InvalidArgumentError(
			ctx,
			"invalid format user id",
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid format user id",
			})
	}

	InitiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)

		return nil, problems.UnauthenticatedError(ctx, "initiator ID format is invalid")
	}

	err = s.app.AdminDeleteUserSessions(ctx, InitiatorID, userId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete sessions for user %s", req.UserId)

		return nil, err
	}

	logger.Log(ctx).Warnf("User sessions deleted by admin %s", InitiatorID)

	return &emptypb.Empty{}, nil
}
