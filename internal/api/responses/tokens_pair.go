package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func TokensPair(pair models.TokensPair) *svc.TokensPair {
	return &svc.TokensPair{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
