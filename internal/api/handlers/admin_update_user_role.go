package handlers

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) AdminUpdateUserRole(ctx context.Context, req *svc.AdminUpdateUserRoleRequest) (*svc.UserResponse, error) {
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

	user, err := s.app.AdminUpdateUserRole(ctx, meta.InitiatorID, userId, role)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Warnf("user %s role updated to %s by %s", userId, role, meta.InitiatorID)
	return responses.User(user), nil
}
