package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPair, error) {
	curToken := req.RefreshToken

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.Id)

		return nil, responses.UnauthenticatedError(
			ctx,
			"invalid initiator ID format",
		)
	}

	sessionID, err := uuid.Parse(req.Initiator.Session)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid session ID format: %s", req.Initiator.Session)

		return nil, responses.UnauthenticatedError(
			ctx,
			"invalid session ID format",
		)
	}

	session, tokensPair, err := s.app.Refresh(ctx, initiatorID, sessionID, req.Agent, curToken)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to refresh session token")

		return nil, responses.AppError(ctx, err)
	}

	logger.Log(ctx).Infof("Session %s refreshed successfully", session.ID)
	return responses.TokensPair(tokensPair), nil
}
