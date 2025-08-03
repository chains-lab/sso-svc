package responses

import (
	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func TokensPair(pair models.TokensPair) *svc.TokensPair {
	return &svc.TokensPair{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
