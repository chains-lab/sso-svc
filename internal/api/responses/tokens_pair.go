package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/auth"
)

func TokensPair(pair models.TokensPair) *auth.TokensPairResponse {
	return &auth.TokensPairResponse{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
