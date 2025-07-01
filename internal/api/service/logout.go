package service

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/responses"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	err := s.app.DeleteSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("failed to delete session %s for user %s", meta.SessionID, meta.InitiatorID)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("User %s Session %s deleted successfully", meta.InitiatorID, meta.SessionID)
	return &emptypb.Empty{}, nil
}
