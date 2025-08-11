package admin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/admin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/guard"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteSessionsByAdmin(ctx context.Context, req *svc.DeleteSessionsByAdminRequest) (*emptypb.Empty, error) {
	initiatorID, err := guard.AllowedRoles(ctx, req.Initiator, "delete user sessions by admin",
		roles.Admin, roles.SuperUser)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	err = s.app.DeleteUserSessions(ctx, user.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("failed to delete sessions for user %s", req.UserId)

		return nil, err
	}

	logger.Log(ctx).Warnf("User sessions deleted by admin %s", initiatorID)

	return &emptypb.Empty{}, nil
}
