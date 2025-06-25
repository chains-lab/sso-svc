package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminDeleteUserSession(ctx context.Context, req *auth.AdminDeleteUserSessionRequest) (*auth.Empty, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		Log(ctx, requestID).WithError(err).Errorf("invalid format user id: %s", req.UserId)

		return nil, status.Errorf(codes.InvalidArgument, "invalid format user id: %s", req.UserId)
	}

	err = a.app.DeleteSession(ctx, meta.InitiatorID, userId)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to delete session")

		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).WithField("user_id", userId).Infof("User sessions deleted by admin %s", meta.InitiatorID)
	return &auth.Empty{}, nil
}
