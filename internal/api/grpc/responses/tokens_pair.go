package responses

import (
	sessionProto "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func TokensPair(pair models.TokensPair) *sessionProto.TokensPair {
	return &sessionProto.TokensPair{
		AccessToken:  pair.Access,
		RefreshToken: pair.Refresh,
	}
}
