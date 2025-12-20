package responses

import (
	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/resources"
)

func TokensPair(m entity.TokensPair) resources.TokensPair {
	resp := resources.TokensPair{
		Data: resources.TokensPairData{
			Id:   m.SessionID,
			Type: resources.TokensPairType,
			Attributes: resources.TokensPairDataAttributes{
				AccessToken:  m.Access,
				RefreshToken: m.Refresh,
			},
		},
	}

	return resp
}
