package response

import (
	userProto "github.com/chains-lab/sso-proto/gen/go/user"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func User(user models.User) *userProto.User {
	return &userProto.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		Verified:  user.Verified,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
