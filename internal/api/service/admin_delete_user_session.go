package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) AdminDeleteUserSession(ctx context.Context, req *svc.AdminDeleteUserSessionRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("invalid format user id: %s", req.UserId)

		return &emptypb.Empty{}, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	err = s.app.DeleteSession(ctx, meta.InitiatorID, userId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Error("failed to delete session")

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).WithField("user_id", userId).Infof("User sessions deleted by admin %s", meta.InitiatorID)
	return &emptypb.Empty{}, nil
}
