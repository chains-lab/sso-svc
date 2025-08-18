package admin

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetUser(ctx context.Context, req *svc.GetUserRequest) (*userProto.User, error) {
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

	user, err := s.app.AdminGetUser(ctx, initiator.ID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error fetching user with ID %s", userID)

		return nil, problems.InternalError(ctx)
	}

	return response.User(user), nil
}
