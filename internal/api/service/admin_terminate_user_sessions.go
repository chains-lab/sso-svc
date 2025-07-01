package service

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) TerminateUserSessionsByAdmin(ctx context.Context, req *svc.TerminateUserSessionsByAdminRequest) (*emptypb.Empty, error) {
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, responses.BadRequestError(ctx, meta.RequestID, responses.Violation{
			Field:       "user_id",
			Description: "invalid format user id",
		})
	}

	err = s.app.AdminTerminateSessions(ctx, meta.InitiatorID, userId)
	if err != nil {
		Log(ctx, meta.RequestID).WithError(err).Errorf("failed to terminate sessions for user %s", req.UserId)

		return nil, responses.AppError(ctx, meta.RequestID, err)
	}

	Log(ctx, meta.RequestID).Warnf("User sessions terminated by admin %s", meta.InitiatorID)

	return &emptypb.Empty{}, nil
}
