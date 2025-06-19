package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) Logout(ctx context.Context, req *sso.SessionRequest) (*sso.Empty, error) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session id")
	}

	err = h.app.DeleteSession(ctx, userID, sessionID)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	log.Infof("User %s Session %s deleted successfully", userID, sessionID)
	return &sso.Empty{}, nil
}
