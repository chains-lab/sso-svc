package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/sso"
)

func User(user models.User) *sso.UserResponse {
	return &sso.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		Verified:     user.Verified,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}
}
