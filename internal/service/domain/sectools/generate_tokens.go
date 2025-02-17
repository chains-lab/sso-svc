package sectools

import (
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
)

func GenerateUserPairTokens(svc *service.Service, account *models.Account, deviceID *string) (*string, *string, error) {
	tokenAccess, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		account.ID.String(),
		svc.Config.JWT.AccessToken.TokenLifetime,
		nil,
		&account.Role,
		deviceID,
	)
	if err != nil {
		return nil, nil, err
	}

	tokenRefresh, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		account.ID.String(),
		svc.Config.JWT.RefreshToken.TokenLifetime,
		nil,
		&account.Role,
		deviceID,
	)
	if err != nil {
		return nil, nil, err
	}

	return &tokenAccess, &tokenRefresh, nil
}
