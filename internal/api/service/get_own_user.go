package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetUser(ctx context.Context, _ *emptypb.Empty) (*svc.User, error) {
	meta := Meta(ctx)

	user, err := s.app.GetUserByID(ctx, meta.InitiatorID)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Error("failed to get user by ID")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	return responses.User(user), nil
}
