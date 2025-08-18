package user

import (
	"context"

	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetOwnUserData(ctx context.Context, _ *emptypb.Empty) (*userProto.User, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	resp, err := s.app.GetUserByID(ctx, initiator.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user by ID")

		return nil, err
	}

	return response.User(resp), nil
}
