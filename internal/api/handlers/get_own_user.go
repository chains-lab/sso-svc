package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
)

func (a Service) GetUser(ctx context.Context, req *svc.Empty) (*svc.UserResponse, error) {
	meta := Meta(ctx)

	user, err := a.app.GetUserByID(ctx, meta.InitiatorID)
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
