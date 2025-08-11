package admin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/admin"
	userProto "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/guard"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
)

func (s Service) UpdateUserVerified(ctx context.Context, req *svc.UpdateUserVerifiedRequest) (*userProto.User, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "update user verified",
		roles.Admin, roles.SuperUser, roles.Moder)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	user, err = s.app.UpdateUserVerified(ctx, user.ID, req.Verified)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user verified status")

		return nil, err
	}

	logger.Log(ctx).Warnf("user %s verified status updated to %v successfully", user.ID, req.Verified)

	return response.User(user), nil
}
