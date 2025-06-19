package handlers

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) AdminUpdateUserRole(ctx context.Context, req *sso.UpdateUserRoleRequest) (*sso.Empty, error) {
	requestID := uuid.New()

	initiatorID, err := uuid.Parse(req.InitiatorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format initiator id: %v", err))
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format user id: %v", err))
	}

	role, err := roles.ParseRole(req.Role)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid format role: %v", err))
	}

	err = h.app.AdminUpdateUserRole(ctx, initiatorID, userId, role)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	h.log.WithField("request_id", requestID).Warnf("user %s role updated to %s by %s", userId, role, initiatorID)
	return &sso.Empty{}, nil
}
