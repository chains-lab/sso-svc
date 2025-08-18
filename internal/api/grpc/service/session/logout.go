package session

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	err = s.app.DeleteUserSession(ctx, initiator.ID, initiator.SessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete session %s for user %s", initiator.ID, initiator.SessionID)

		return nil, err
	}

	logger.Log(ctx).Infof("User %s Session %s deleted successfully", initiator.ID, initiator.SessionID)

	return &emptypb.Empty{}, nil
}
