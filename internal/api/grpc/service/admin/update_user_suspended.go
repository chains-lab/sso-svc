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

func (s Service) UpdateUserSuspended(ctx context.Context, req *svc.UpdateUserSuspendedRequest) (*userProto.User, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "update user suspended",
		roles.Admin, roles.SuperUser, roles.Moder)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	user, err = s.app.UpdateUserSuspended(ctx, user.ID, req.Suspended)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update user suspended status")

		return nil, err
	}

	logger.Log(ctx).Warnf("user %s suspended status updated to %v successfully", user.ID, req.Suspended)

	return response.User(user), nil
}
