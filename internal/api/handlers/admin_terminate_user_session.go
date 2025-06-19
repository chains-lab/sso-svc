package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminTerminateUserSessions(ctx context.Context, req *sso.AdminUserRequest) (*sso.Empty, error) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %s", req.UserId)
	}

	err = h.app.TerminateSessionsByAdmin(ctx, userId)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	log.WithField("user_id", userId).Info("User sessions terminated by admin")

	return &sso.Empty{}, nil
}
