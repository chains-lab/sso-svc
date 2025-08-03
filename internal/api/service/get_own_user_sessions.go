package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) GetOwnUserSessions(ctx context.Context, req *svc.GetOwnUserSessionsRequest) (*svc.SessionsList, error) {
	meta := Meta(ctx)

	sessions, err := s.app.GetUserSessions(ctx, meta.InitiatorID, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get user sessions")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.SessionList(sessions, req.Pagination.Page), nil
}
