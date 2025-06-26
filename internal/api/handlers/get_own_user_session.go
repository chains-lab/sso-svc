package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (a Service) GetUserSession(ctx context.Context, req *svc.Empty) (*svc.SessionResponse, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	session, err := a.app.GetSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("delete session %s for user %s", meta.SessionID, meta.InitiatorID)

	return responses.Session(session), nil
}
