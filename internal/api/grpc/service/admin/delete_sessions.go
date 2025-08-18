package admin

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteSessions(ctx context.Context, req *svc.DeleteSessionsRequest) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, problems.InvalidArgumentError(ctx, "user_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	err = s.app.DeleteUserSessions(ctx, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete sessions for user %s", req.UserId)

		return nil, err
	}

	logger.Log(ctx).Warnf("User sessions deleted by admin %s", initiator.ID)

	return &emptypb.Empty{}, nil
}
