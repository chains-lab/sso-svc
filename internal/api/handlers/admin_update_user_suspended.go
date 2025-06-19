package handlers

import (
	"context"
	"fmt"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminUpdateUserSuspended(ctx context.Context, req *sso.UpdateUserSuspendedRequest) (*sso.Empty, error) {
	requestID := uuid.New()

	initiatorID, err := uuid.Parse(req.InitiatorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format initiator id: %v", err))
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	err = h.app.AdminUpdateUserSuspended(ctx, initiatorID, userID, req.Suspended)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Warnf("user %s suspended status updated to %t by %s", userID, req.Suspended, initiatorID)
	return &sso.Empty{}, nil
}
