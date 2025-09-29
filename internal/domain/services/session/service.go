package session

import (
	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/google/uuid"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	ParseRefreshClaims(enc string) (auth.UsersClaims, error)

	GenerateAccess(
		userID uuid.UUID,
		sessionID uuid.UUID,
		idn string,
	) (string, error)

	GenerateRefresh(
		userID uuid.UUID,
		sessionID uuid.UUID,
		role string,
	) (string, error)
}

type Service struct {
	jwt JWTManager
	db  data.Database
}

func New(db data.Database, jwt JWTManager) Service {
	return Service{
		jwt: jwt,
		db:  db,
	}
}
