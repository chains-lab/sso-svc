package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminTerminateUserSessions(ctx context.Context, req *sso.TerminateUserSessionByAdminRequest) (*sso.Empty, error) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format user id: %s", req.UserId)
	}

	initiatorId, err := uuid.Parse(req.InitiatorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format initiator id: %s", req.InitiatorId)
	}

	err = h.app.AdminTerminateSessions(ctx, initiatorId, userId)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	log.WithField("user_id", userId).Warnf("User sessions terminated by admin %s", initiatorId)

	return &sso.Empty{}, nil
}
