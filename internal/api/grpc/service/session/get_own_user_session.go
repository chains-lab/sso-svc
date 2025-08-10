package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetOwnSession(ctx context.Context, req *svc.GetOwnSessionRequest) (*svc.Session, error) {
	InitiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.Id)
	}

	SessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid session ID format: %s", req.SessionId)

		return nil, responses.InvalidArgumentError(
			ctx,
			"invalid session ID format",
			&errdetails.BadRequest_FieldViolation{
				Field:       "session_id",
				Description: "invalid UUID format for session ID",
			})
	}

	session, err := s.app.GetSession(ctx, InitiatorID, SessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user session")

		return nil, responses.AppError(ctx, err)
	}

	return responses.Session(session), nil
}
