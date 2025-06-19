package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminUpdateUserVerified(ctx context.Context, req *sso.UpdateUserVerifiedRequest) (*sso.Empty, error) {
	requestID := uuid.New()

	initiatorID, err := uuid.Parse(req.InitiatorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid format initiator id: "+err.Error())
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid format user id: "+err.Error())
	}

	err = h.app.AdminUpdateUserVerified(ctx, initiatorID, userID, req.Verified)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Warnf("user %s verified status updated to %t by %s", userID, req.Verified, initiatorID)
	return &sso.Empty{}, nil
}
