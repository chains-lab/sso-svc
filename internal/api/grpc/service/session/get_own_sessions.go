package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
)

func (s Service) GetOwnSessions(ctx context.Context, req *svc.GetOwnSessionsRequest) (*svc.SessionsList, error) {
	InitiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", req.Initiator.UserId)

		return nil, problem.UnauthenticatedError(ctx, "initiator ID format is invalid")
	}

	session, pag, err := s.app.GetUserSessions(ctx, InitiatorID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user session")

		return nil, err
	}

	return response.SessionList(session, pag), nil
}
