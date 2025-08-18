package admin

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	sesionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/chains-lab/sso-svc/internal/pagination"
)

func (s Service) GetSessions(ctx context.Context, req *svc.GetSessionsRequest) (*sesionProto.SessionsList, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, problems.InvalidArgumentError(ctx, "user_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	sessions, pag, err := s.app.AdminGetUserSessions(ctx, initiator.ID, userID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to get sessions for user %s", req.UserId)

		return nil, err
	}

	return response.SessionList(sessions, pag), nil
}
