package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnSession(ctx context.Context, req *svc.DeleteOwnSessionRequest) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid session ID format")

		return nil, problems.InvalidArgumentError(
			ctx,
			"invalid session ID format",
			&errdetails.BadRequest_FieldViolation{
				Field:       "session_id",
				Description: "invalid UUID format for session ID",
			})
	}

	err = s.app.DeleteUserSession(ctx, initiator.ID, sessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to delete user session")

		return nil, err
	}

	logger.Log(ctx).Infof("delete session %s for user %s", sessionID, initiator.ID)

	return &emptypb.Empty{}, nil
}
