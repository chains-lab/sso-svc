package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (s Service) TerminateUserSessions(ctx context.Context, req *svc.Empty) (*svc.Empty, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	err := s.app.TerminateUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("User sessions terminated for user ID: %s", meta.InitiatorID)

	return &svc.Empty{}, nil
}
