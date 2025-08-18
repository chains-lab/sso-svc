package response

import (
	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func User(user models.User) *userProto.User {
	return &userProto.User{
		Id:            user.ID.String(),
		Role:          user.Role,
		Email:         user.Email,
		EmailVerified: user.EmailVer,
		CreatedAt:     user.CreatedAt.String(),
		UpdatedAt:     user.EmailUpdatedAt.String(),
	}
}
