package jwtmanager

import (
	"github.com/chains-lab/restkit/token"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) EncryptAccess(token string) (string, error) {
	return encryptAESGCM(token, []byte(s.accessSK))
}

func (s Service) GenerateAccess(user models.User, sessionID uuid.UUID) (string, error) {
	return token.GenerateUserJWT(token.GenerateUserJwtRequest{
		Issuer:    s.iss,
		UserID:    user.ID,
		SessionID: sessionID,
		Role:      user.SysRole,
		Ttl:       s.accessTTL,
	}, s.accessSK)
}
