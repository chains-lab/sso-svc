package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPair, error) {
	meta := Meta(ctx)

	curToken := req.RefreshToken

	session, tokensPair, err := s.app.Refresh(ctx, meta.InitiatorID, meta.SessionID, req.Agent, curToken)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to refresh session token")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("Session %s refreshed successfully", session.ID)
	return responses.TokensPair(tokensPair), nil
}
