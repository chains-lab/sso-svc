package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (a Service) Logout(ctx context.Context, req *svc.Empty) (*svc.Empty, error) {
	requestID := uuid.New()
	log := Log(ctx, requestID)
	meta := Meta(ctx)

	err := a.app.DeleteSession(ctx, meta.InitiatorID, meta.SessionID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	log.Infof("User %s Session %s deleted successfully", meta.InitiatorID, meta.SessionID)
	return &svc.Empty{}, nil
}
