package responses

import (
	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/resources"
)

func AccountEmailData(ae entity.AccountEmail) resources.AccountEmail {
	return resources.AccountEmail{
		Data: resources.AccountEmailData{
			Id:   ae.AccountID,
			Type: resources.AccountEmailType,
			Attributes: resources.AccountEmailDataAttributes{
				Email:     ae.Email,
				Verified:  ae.Verified,
				UpdatedAt: ae.UpdatedAt,
			},
		},
	}
}
