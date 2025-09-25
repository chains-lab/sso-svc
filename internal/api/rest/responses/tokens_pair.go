package responses

import (
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/resources"
)

func TokensPair(m models.TokensPair) resources.TokensPair {
	resp := resources.TokensPair{
		Data: resources.TokensPairData{
			Id:   m.SessionID.String(),
			Type: resources.TokensPairType,
			Attributes: resources.TokensPairDataAttributes{
				AccessToken:  m.Access,
				RefreshToken: m.Refresh,
			},
		},
	}

	return resp
}
