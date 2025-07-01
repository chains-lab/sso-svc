package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"

	"github.com/google/uuid"
)

func (s Service) AdminGetUserSession(ctx context.Context, req *svc.AdminGetUserSessionRequest) (*svc.Session, error) {
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return &svc.Session{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return &svc.Session{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "session_id",
			Description: "invalid format session id",
		})
	}

	session, err := s.app.GetSession(ctx, userId, sessionId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("failed to retrieve session for user %s", userId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Infof("Retrieved session for user %s by admin %s", userId, meta.InitiatorID)

	return responses.Session(session), nil
}
