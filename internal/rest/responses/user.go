package responses

import (
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/resources"
)

func User(m models.User) resources.User {
	resp := resources.User{
		Data: resources.UserData{
			Id:   m.ID,
			Type: resources.UserTepe,
			Attributes: resources.UserDataAttributes{
				Email:       m.Email,
				Role:        m.SysRole,
				CreatedAt:   m.CreatedAt,
				CityId:      m.CityID,
				CityRole:    m.CityRole,
				CompanyId:   m.CompanyID,
				CompanyRole: m.CompanyRole,
				UpdatedAt:   m.UpdatedAt,
			},
		},
	}

	return resp
}
