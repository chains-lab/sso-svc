package service

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/app/ape"
)

func (s Service) CreateAdminByAdmin(ctx context.Context, req *svc.CreateAdminByAdminRequest) (*svc.User, error) {
	meta := Meta(ctx)
	if meta.Role != roles.SuperUser {
		return nil, responses.AppError(ctx, meta.RequestID, ape.ErrorNoPermission(
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

	user, err := s.app.CreateUser(ctx, req.Email, role)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to create admin user")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Warnf("admin user %s created successfully", user.ID)
	return responses.User(user), nil
}
