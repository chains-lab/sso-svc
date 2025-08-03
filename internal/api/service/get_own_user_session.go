package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetUserSession(ctx context.Context, _ *emptypb.Empty) (*svc.Session, error) {
	meta := Meta(ctx)

	session, err := s.app.GetSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get user session")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.Session(session), nil
}
