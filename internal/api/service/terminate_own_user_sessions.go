package service

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/responses"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) TerminateUserSessions(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.TerminateUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to terminate user sessions")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("User sessions terminated for user ID: %s", meta.InitiatorID)
	return &emptypb.Empty{}, nil
}
