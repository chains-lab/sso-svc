package handlers

import (
	"context"

	"github.com/chains-lab/proto-storage/gen/go/auth"
)

func (a Service) GetUser(ctx context.Context, req *auth.Empty) (*auth.UserResponse, error) {
	meta := Meta(ctx)

	user, err := a.app.GetUserByID(ctx, meta.InitiatorID)
	if err != nil {
		return nil, err
	}

	return &auth.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}
