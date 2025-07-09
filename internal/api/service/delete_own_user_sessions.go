package service

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnUserSessions(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.DeleteUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to delete user sessions")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("User sessions deleted for user ID: %s", meta.InitiatorID)
	return &emptypb.Empty{}, nil
}
