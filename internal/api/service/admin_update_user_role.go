package service

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (s Service) AdminUpdateUserRole(ctx context.Context, req *svc.AdminUpdateUserRoleRequest) (*svc.User, error) {
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return &svc.User{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	role, err := roles.ParseRole(req.Role)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("invalid role format: %s", req.Role)

		return &svc.User{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "role",
			Description: "invalid format role",
		})
	}

	user, err := s.app.AdminUpdateUserRole(ctx, meta.InitiatorID, userId, role)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("failed to update user role for user %s", userId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Warnf("user %s role updated to %s by %s", userId, role, meta.InitiatorID)
	return responses.User(user), nil
}
