package user

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problem"
	"github.com/chains-lab/sso-svc/internal/api/grpc/response"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetOwnUserData(ctx context.Context, req *svc.GetOwnUserDataRequest) (*svc.User, error) {
	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, problem.UnauthenticatedError(
			ctx,
			"invalid initiator ID format",
		)
	}

	user, err := s.app.GetUserByID(ctx, initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user by ID")

		return nil, err
	}

	return response.User(user), nil
}
