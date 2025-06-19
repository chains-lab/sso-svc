package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) GetUserSession(ctx context.Context, req *sso.SessionRequest) (*sso.SessionResponse, error) {
	requestID := uuid.New()

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session id: %s", req.SessionId)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %s", req.UserId)
	}

	session, err := h.app.GetUserSession(ctx, userID, sessionID)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Infof("delete session %s for user %s", sessionID, userID)

	return responses.Session(session), nil
}
