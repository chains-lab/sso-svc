package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/sso"
)

func TokensPair(pair models.TokensPair) *sso.TokensPairResponse {
	return &sso.TokensPairResponse{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
