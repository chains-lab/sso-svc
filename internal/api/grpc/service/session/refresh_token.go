package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPair, error) {
	curToken := req.RefreshToken

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)

		return nil, problem.UnauthenticatedError(ctx, "invalid initiator ID format")
	}

	sessionID, err := uuid.Parse(req.Initiator.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid session ID format: %s", req.Initiator.SessionId)

		return nil, problem.UnauthenticatedError(ctx, "invalid session ID format")
	}

	session, tokensPair, err := s.app.Refresh(ctx, initiatorID, sessionID, req.Agent, curToken)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to refresh session token")

		return nil, err
	}

	logger.Log(ctx).Infof("Session %s refreshed successfully", session.ID)
	return response.TokensPair(tokensPair), nil
}
