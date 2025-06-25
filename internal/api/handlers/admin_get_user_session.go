package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminGetUserSession(ctx context.Context, req *auth.AdminGetUserSessionRequest) (*auth.SessionResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format user id: %s", req.UserId)
	}

	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format session id: %s", req.SessionId)
	}

	session, err := a.app.GetSession(ctx, userId, sessionId)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("Retrieved session for user %s by admin %s", userId, meta.InitiatorID)

	return responses.Session(session), nil
}
