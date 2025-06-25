package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"

	"github.com/google/uuid"
)

func (a Service) GetUserSessions(ctx context.Context, req *auth.Empty) (*auth.SessionsListResponse, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	sessions, err := a.app.SelectUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.SessionList(sessions), nil
}
