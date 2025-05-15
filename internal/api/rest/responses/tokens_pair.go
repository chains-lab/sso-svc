package responses

import "github.com/chains-lab/chains-auth/resources"

func TokensPair(access, refresh string) resources.TokensPair {
	return resources.TokensPair{
		Data: resources.TokensPairData{
			Type: resources.TokensPairType,
			Attributes: resources.TokensPairDataAttributes{
				AccessToken:  access,
				RefreshToken: refresh,
			},
		},
	}
}
