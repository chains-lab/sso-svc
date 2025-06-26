package handlers

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminUpdateUserSubscription(ctx context.Context, req *svc.AdminUpdateUserSubscriptionRequest) (*svc.UserResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	subscriptionID, err := uuid.Parse(req.Subscription)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format subscription id: %v", err))
	}

	user, err := a.app.AdminUpdateUserSubscription(ctx, meta.InitiatorID, userID, subscriptionID)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("user %s subscription updated to %s by %s", userID, req.Subscription, meta.InitiatorID)
	return responses.User(user), nil
}
