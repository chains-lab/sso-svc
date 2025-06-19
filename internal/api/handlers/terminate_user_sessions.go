package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) TerminateUserSessions(ctx context.Context, req *sso.UserRequest) (*sso.Empty, error) {
	requestID := uuid.New()

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %s", req.UserId)
	}

	err = h.app.TerminateUserSessions(ctx, userID)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Infof("User sessions terminated for user ID: %s", userID)

	return &sso.Empty{}, nil
}
