package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
)

func (s Service) AdminGetUser(ctx context.Context, req *svc.AdminGetUserRequest) (*svc.UserResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	user, err := s.app.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &svc.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}
