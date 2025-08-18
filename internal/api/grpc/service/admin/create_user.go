package admin

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/svc/admin"
	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) CreateUser(ctx context.Context, req *svc.CreateUserRequest) (*userProto.User, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	user, err := s.app.AdminCreateUser(ctx, initiator.ID, app.AdminCreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Verified: true,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to create user")

		return nil, err
	}

	logger.Log(ctx).Infof("user %s created successfully", user.ID)

	return response.User(user), nil
}
