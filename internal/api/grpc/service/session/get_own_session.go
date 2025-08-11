package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetOwnSession(ctx context.Context, req *svc.GetOwnSessionRequest) (*svc.Session, error) {
	InitiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)

		return nil, problems.UnauthenticatedError(ctx, "initiator ID format is invalid")
	}

	SessionID, err := uuid.Parse(req.Initiator.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid session ID format: %s", req.Initiator.SessionId)

		return nil, problems.UnauthenticatedError(ctx, "session ID format is invalid")
	}

	session, err := s.app.GetUserSession(ctx, InitiatorID, SessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user session")

		return nil, err
	}

	return responses.Session(session), nil
}
