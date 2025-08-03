package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteOwnUserSession(ctx context.Context, req *svc.DeleteOwnUserSessionRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("invalid session ID format")

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "session_id",
			Description: "invalid UUID format for session ID",
		})
	}

	err = s.app.DeleteSession(ctx, meta.InitiatorID, sessionID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to delete user session")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Infof("delete session %s for user %s", meta.SessionID, meta.InitiatorID)

	return &emptypb.Empty{}, nil
}
