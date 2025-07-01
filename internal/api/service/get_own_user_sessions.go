package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetUserSessions(ctx context.Context, _ *emptypb.Empty) (*svc.SessionsList, error) {
	meta := Meta(ctx)

	sessions, err := s.app.SelectUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get user sessions")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.SessionList(sessions), nil
}
