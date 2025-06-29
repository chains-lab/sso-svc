package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"

	"github.com/google/uuid"
)

func (s Service) GetUserSessions(ctx context.Context, req *svc.Empty) (*svc.SessionsListResponse, error) {
	requestID := uuid.New()

	meta := Meta(ctx)

	sessions, err := s.app.SelectUserSessions(ctx, meta.InitiatorID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.SessionList(sessions), nil
}
