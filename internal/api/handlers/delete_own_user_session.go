package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (a Service) DeleteUserSession(ctx context.Context, req *svc.Empty) (*svc.SessionsListResponse, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	err := a.app.DeleteSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	sessions, err := a.app.SelectUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("delete session %s for user %s", meta.SessionID, meta.InitiatorID)

	return responses.SessionList(sessions), nil
}
