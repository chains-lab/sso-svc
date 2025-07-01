package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPair, error) {
	meta := Meta(ctx)

	curToken := req.RefreshToken

	session, tokensPair, err := s.app.Refresh(ctx, meta.InitiatorID, meta.SessionID, req.Agent, curToken)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to refresh session token")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("Session %s refreshed successfully", session.ID)
	return responses.TokensPair(tokensPair), nil
}
