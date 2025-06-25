package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminTerminateUserSessions(ctx context.Context, req *auth.AdminTerminateUserSessionsRequest) (*auth.Empty, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format user id: %s", req.UserId)
	}

	err = a.app.AdminTerminateSessions(ctx, meta.InitiatorID, userId)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("User sessions terminated by admin %s", meta.InitiatorID)

	return &auth.Empty{}, nil
}
