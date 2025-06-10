package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/resources"
)

func User(user models.User) resources.User {
	return resources.User{
		Data: resources.UserData{
			Id:   user.ID.String(),
			Type: resources.UserType,
			Attributes: resources.UserDataAttributes{
				Email:     user.Email,
				Role:      string(user.Role),
				UpdatedAt: user.UpdatedAt,
				CreatedAt: user.CreatedAt,
			},
		},
	}
}
