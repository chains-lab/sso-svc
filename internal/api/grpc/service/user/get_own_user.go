package user

import (
	"context"

	svc "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetOwnUser(ctx context.Context, req *svc.GetOwnUserDataRequest) (*svc.User, error) {
	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to parse initiator ID")

		return nil, responses.AppError(ctx, err)
	}

	user, err := s.app.GetUserByID(ctx, initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get user by ID")

		return nil, responses.AppError(ctx, err)
	}

	return responses.User(user), nil
}
