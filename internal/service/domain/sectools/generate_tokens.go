package sectools

import (
	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
)

func GenerateTokens(service service.Service, account models.Account, deviceID uuid.UUID) (tokenAccess string, tokenRefresh string, err error) {
	tokenAccess, err = service.TokenManager.GenerateJWT(account.ID, deviceID, account.Role, service.Config.JWT.AccessToken.TokenLifetime, service.Config.JWT.AccessToken.SecretKey)
	if err != nil {
		return "", "", err
	}

	tokenRefresh, err = service.TokenManager.GenerateJWT(account.ID, deviceID, account.Role, service.Config.JWT.RefreshToken.TokenLifetime, service.Config.JWT.RefreshToken.SecretKey)
	if err != nil {
		return "", "", err
	}

	return tokenAccess, tokenRefresh, nil
}
