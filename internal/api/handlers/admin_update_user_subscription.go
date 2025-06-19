package handlers

import (
	"context"
	"fmt"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminUpdateUserSubscription(ctx context.Context, req *sso.UpdateUserSubscriptionRequest) (*sso.Empty, error) {
	requestID := uuid.New()

	initiatorID, err := uuid.Parse(req.InitiatorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format initiator id: %v", err))
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	subscriptionID, err := uuid.Parse(req.Subscription)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format subscription id: %v", err))
	}

	err = h.app.AdminUpdateUserSubscription(ctx, initiatorID, userID, subscriptionID)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Warnf("user %s subscription updated to %s by %s", userID, req.Subscription, initiatorID)
	return &sso.Empty{}, nil
}
