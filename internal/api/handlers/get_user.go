package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) GetUser(ctx context.Context, req *sso.UserRequest) (*sso.UserResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	user, err := h.app.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &sso.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}
