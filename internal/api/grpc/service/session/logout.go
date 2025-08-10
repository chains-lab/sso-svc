package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) Logout(ctx context.Context, req *svc.LogoutRequest) (*emptypb.Empty, error) {
	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)

		return nil, problems.UnauthenticatedError(
			ctx,
			"invalid initiator ID format",
		)
	}

	sessionID, err := uuid.Parse(req.Initiator.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid session ID format: %s", req.Initiator.SessionId)

		return nil, problems.UnauthenticatedError(
			ctx,
			"invalid session ID format",
		)
	}

	err = s.app.DeleteUserSession(ctx, initiatorID, sessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete session %s for user %s", sessionID, initiatorID)

		return nil, problems.AppError(ctx, err)
	}

	logger.Log(ctx).Infof("User %s Session %s deleted successfully", initiatorID, sessionID)
	return &emptypb.Empty{}, nil
}
