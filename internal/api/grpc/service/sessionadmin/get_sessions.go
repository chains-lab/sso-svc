package sessionadmin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	sesionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	svc "github.com/chains-lab/sso-proto/gen/go/svc/sessionadmin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/guard"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"

	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
)

func (s Service) GetSessions(ctx context.Context, req *svc.GetSessionsRequest) (*sesionProto.SessionsList, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "get user sessions by admin",
		roles.Admin, roles.SuperUser)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	sessions, pag, err := s.app.GetUserSessions(ctx, user.ID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to get sessions for user %s", req.UserId)

		return nil, err
	}

	return response.SessionList(sessions, pag), nil
}
