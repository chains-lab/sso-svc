package jwtmanager

import (
	"fmt"

	"github.com/chains-lab/restkit/token"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GenerateRefresh(user models.User, sessionID uuid.UUID) (string, error) {
	return token.GenerateUserJWT(token.GenerateUserJwtRequest{
		Issuer:   s.iss,
		Audience: []string{""}, //TODO: need be om;y api-gateway while we have dev stage this field must be empty
		Role:     user.SysRole,
		Ttl:      s.refreshTTL,
	}, s.refreshSK)
}

func (s Service) EncryptRefresh(token string) (string, error) {
	return encryptAESGCM(token, []byte(s.refreshSK))
}

func (s Service) DecryptRefresh(encryptedToken string) (string, error) {
	raw, err := decryptAESGCM(encryptedToken, []byte(s.refreshSK))
	if err != nil {
		return "", fmt.Errorf("decrypt refresh: %w", err)
	}

	return raw, nil
}

func (s Service) ParseRefreshClaims(tokenStr string) (token.UsersClaims, error) {
	return token.VerifyUserJWT(tokenStr, s.refreshSK)
}
