package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminGetUserSessions(ctx context.Context, req *auth.AdminGetUserSessionsRequest) (*auth.SessionsListResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	sessions, err := a.app.SelectUserSessions(ctx, userID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("retrieved sessions for user %s by admin %s", req.UserId, meta.InitiatorID)

	return responses.SessionList(sessions), nil
}
