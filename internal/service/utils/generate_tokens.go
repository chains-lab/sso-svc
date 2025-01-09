package utils

import (
	"github.com/cifra-city/sso-oauth/internal/config"
	"github.com/cifra-city/sso-oauth/internal/data/sql/repositories/sqlcore"
	"github.com/google/uuid"
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
