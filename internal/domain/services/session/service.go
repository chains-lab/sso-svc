package session

import (
	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/chains-lab/sso-svc/internal/domain/services/session/jwtmanager"
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
		emailVerified bool,
	) (string, error)
}

type Service struct {
	db  data.Database
	jwt JWTManager
}

func NewService(db data.Database, manager jwtmanager.Manager) Service {
	return Service{
		db:  db,
		jwt: manager,
	}
}
