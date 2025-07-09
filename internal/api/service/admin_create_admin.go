package service

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/ape"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) CreateUserByAdmin(ctx context.Context, req *svc.CreateUserByAdminRequest) (*svc.User, error) {
	meta := Meta(ctx)
	if meta.Role != roles.SuperUser {
		return nil, responses.AppError(ctx, meta.RequestID, ape.RaiseNoPermission(
			fmt.Errorf("only superuser can create admin user, current role: %s", meta.Role)),
		)
	}

	role, err := roles.ParseRole(req.Role)
	if err != nil {
		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "role",
			Description: "invalid role",
		})
	}

	user, err := s.app.AdminCreateUser(ctx, req.Email, role)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to create admin user")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Warnf("admin user %s created successfully", user.ID)
	return responses.User(user), nil
}
