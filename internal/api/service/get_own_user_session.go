package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetUserSession(ctx context.Context, _ *emptypb.Empty) (*svc.Session, error) {
	meta := Meta(ctx)

	session, err := s.app.GetSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to get user session")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Session(session), nil
}
