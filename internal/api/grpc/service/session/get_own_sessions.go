package session

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
)

func (s Service) GetOwnSessions(ctx context.Context, req *svc.GetOwnSessionsRequest) (*svc.SessionsList, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	session, pag, err := s.app.GetUserSessions(ctx, initiator.ID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user session")

		return nil, err
	}

	return response.SessionList(session, pag), nil
}
