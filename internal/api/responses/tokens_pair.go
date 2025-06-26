package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func TokensPair(pair models.TokensPair) *svc.TokensPairResponse {
	return &svc.TokensPairResponse{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
