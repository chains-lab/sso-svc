package token

import (
	"github.com/chains-lab/restkit/token"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/google/uuid"
)

func (s Service) EncryptAccess(token string) (string, error) {
	return encryptAESGCM(token, []byte(s.accessSK))
}

func (s Service) GenerateAccess(user entity.Account, sessionID uuid.UUID) (string, error) {
	return token.GenerateAccountJWT(token.GenerateAccountJwtRequest{
		Issuer:    s.iss,
		AccountID: user.ID,
		//Audience:  []string{"gateway"},
		SessionID: sessionID,
		Role:      user.Role,
		Username:  user.Username,
		Ttl:       s.accessTTL,
	}, s.accessSK)
}
