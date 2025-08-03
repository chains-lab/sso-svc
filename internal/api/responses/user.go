package responses

import (
	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func User(user models.User) *svc.User {
	return &svc.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		Verified:  user.Verified,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
