package admin

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/sso-proto/gen/go/admin"
	userProto "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/guard"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
)

func (s Service) GetUserByAdmin(ctx context.Context, req *svc.GetUserByAdminRequest) (*userProto.User, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "get user by admin by admin",
		roles.Admin, roles.SuperUser)
	if err != nil {
		return nil, err
	}

	_, user, err := s.ComparisonRightsForAdmins(ctx, req.Initiator.UserId, req.UserId)
	if err != nil {
		return nil, err
	}

	return response.User(user), nil
}
