package session

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetSessionsByAdmin(ctx context.Context, req *svc.GetSessionsByAdminRequest) (*svc.SessionsList, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access: only admin or super admin can create user")

		return nil, problems.PermissionDeniedError(ctx, "only admin or super admin can get user sessions")
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, problems.InvalidArgumentError(
			ctx,
			"invalid format user id",
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid format user id",
			})
	}

	sessions, pag, err := s.app.GetUserSessions(ctx, userId, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to get sessions for user %s", req.UserId)

		return nil, err
	}

	return responses.SessionList(sessions, pag), nil
}
