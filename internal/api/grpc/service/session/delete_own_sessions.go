package session

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"

	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnSessions(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	err = s.app.DeleteUserSessions(ctx, initiator.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete sessions for user %s", initiator.ID)

		return nil, err
	}

	logger.Log(ctx).Warnf("User sessions deleted by admin %s", initiator.ID)

	return &emptypb.Empty{}, nil
}
