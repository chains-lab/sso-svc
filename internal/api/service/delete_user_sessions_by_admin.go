package service

import (
	"context"

	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteUserSessionsByAdmin(ctx context.Context, req *svc.DeleteUserSessionsByAdminRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	err = s.app.AdminDeleteSessions(ctx, meta.InitiatorID, userId)
	if err != nil {
		logger.Log(ctx, meta.RequestID).WithError(err).Errorf("failed to delete sessions for user %s", req.UserId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	logger.Log(ctx, meta.RequestID).Warnf("User sessions deleted by admin %s", meta.InitiatorID)

	return &emptypb.Empty{}, nil
}
