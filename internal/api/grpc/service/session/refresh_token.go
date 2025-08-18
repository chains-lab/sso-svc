package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPair, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	curToken := req.RefreshToken

	session, tokensPair, err := s.app.Refresh(
		ctx,
		initiator.ID,
		initiator.SessionID,
		req.Client,
		req.Ip,
		curToken,
	)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to refresh session token")

		return nil, err
	}

	logger.Log(ctx).Infof("Session %s refreshed successfully", session.ID)
	return response.TokensPair(tokensPair), nil
}
