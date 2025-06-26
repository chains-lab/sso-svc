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

func (a Service) AdminUpdateUserSuspended(ctx context.Context, req *svc.AdminUpdateUserSuspendedRequest) (*svc.UserResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	user, err := a.app.AdminUpdateUserSuspended(ctx, meta.InitiatorID, userID, req.Suspended)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("user %s suspended status updated to %t by %s", userID, req.Suspended, meta.InitiatorID)
	return responses.User(user), nil
}
