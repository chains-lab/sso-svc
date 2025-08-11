package user

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetUserByAdmin(ctx context.Context, req *svc.GetUserByAdminRequest) (*svc.User, error) {
	if req.Initiator.Role == roles.Admin || req.Initiator.Role == roles.SuperUser {
		logger.Log(ctx).Error("unauthorized access: only admin or super admin can get user by ID")

		return nil, problems.PermissionDeniedError(ctx, "only admins roles can get user by ID")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse user ID")

		return &svc.User{}, problems.InvalidArgumentError(
			ctx,
			"invalid user ID format",
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "user ID must be a valid UUID",
			})
	}

	user, err := s.app.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user")

		return nil, err
	}

	return &svc.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
