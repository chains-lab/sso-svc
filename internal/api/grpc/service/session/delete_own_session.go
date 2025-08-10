package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnSession(ctx context.Context, req *svc.DeleteOwnSessionRequest) (*emptypb.Empty, error) {

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid session ID format")

		return nil, responses.InvalidArgumentError(
			ctx,
			"invalid session ID format",
			&errdetails.BadRequest_FieldViolation{
				Field:       "session_id",
				Description: "invalid UUID format for session ID",
			})
	}

	InitiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, responses.AppError(ctx, responses.UnauthenticatedError(ctx, "invalid format initiator ID"))
	}

	err = s.app.DeleteSession(ctx, InitiatorID, sessionID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to delete user session")

		return nil, responses.AppError(ctx, err)
	}

	logger.Log(ctx).Infof("delete session %s for user %s", sessionID, InitiatorID)

	return &emptypb.Empty{}, nil
}
