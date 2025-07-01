package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteUserSession(ctx context.Context, _ *emptypb.Empty) (*svc.SessionsList, error) {
	meta := Meta(ctx)

	err := s.app.DeleteSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to delete user session")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	sessions, err := s.app.SelectUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get user sessions")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("delete session %s for user %s", meta.SessionID, meta.InitiatorID)

	return responses.SessionList(sessions), nil
}
