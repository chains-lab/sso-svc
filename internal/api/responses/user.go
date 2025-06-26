package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func User(user models.User) *svc.UserResponse {
	return &svc.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Role:         string(user.Role),
		Subscription: user.Subscription.String(),
		Verified:     user.Verified,
		Suspended:    user.Suspended,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}
}
