package handlers

import (
	"context"
	"fmt"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Service) AdminUpdateUserRole(ctx context.Context, req *auth.AdminUpdateUserRoleRequest) (*auth.UserResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	role, err := roles.ParseRole(req.Role)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format role: %v", err))
	}

	user, err := a.app.AdminUpdateUserRole(ctx, meta.InitiatorID, userId, role)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("user %s role updated to %s by %s", userId, role, meta.InitiatorID)
	return responses.User(user), nil
}
