package utils

import (
	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/data/sql/repositories/sqlcore"
)

func GenerateTokens(service config.Server, account sqlcore.Account, deviceID uuid.UUID) (tokenAccess string, tokenRefresh string, err error) {
	tokenAccess, err = service.TokenManager.GenerateJWT(account.ID, deviceID, string(account.Role), service.Config.JWT.AccessToken.TokenLifetime, service.Config.JWT.AccessToken.SecretKey)
	if err != nil {
		return "", "", err
	}

	tokenRefresh, err = service.TokenManager.GenerateJWT(account.ID, deviceID, string(account.Role), service.Config.JWT.RefreshToken.TokenLifetime, service.Config.JWT.RefreshToken.SecretKey)
	if err != nil {
		return "", "", err
	}

	return tokenAccess, tokenRefresh, nil
}
