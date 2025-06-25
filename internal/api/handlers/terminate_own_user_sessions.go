package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
)

func (a Service) TerminateUserSessions(ctx context.Context, req *auth.Empty) (*auth.Empty, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	err := a.app.TerminateUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("User sessions terminated for user ID: %s", meta.InitiatorID)

	return &auth.Empty{}, nil
}
