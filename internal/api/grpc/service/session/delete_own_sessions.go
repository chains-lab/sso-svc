package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnSessions(ctx context.Context, req *svc.DeleteOwnSessionsRequest) (*emptypb.Empty, error) {
	InitiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)
	}

	err = s.app.DeleteUserSessions(ctx, InitiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete sessions for user %s", InitiatorID)

		return nil, problems.AppError(ctx, err)
	}

	logger.Log(ctx).Warnf("User sessions deleted by admin %s", InitiatorID)

	return &emptypb.Empty{}, nil
}
