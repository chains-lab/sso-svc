package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
)

func (a Service) GetUserSession(ctx context.Context, req *auth.Empty) (*auth.SessionResponse, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	session, err := a.app.GetSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("delete session %s for user %s", meta.SessionID, meta.InitiatorID)

	return responses.Session(session), nil
}
