package responses

import (
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/resources"
)

func User(m models.User) resources.User {
	resp := resources.User{
		Data: resources.UserData{
			Id:   m.ID.String(),
			Type: resources.UserTepe,
			Attributes: resources.UserDataAttributes{
				Email:     m.Email,
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	return resp
}
