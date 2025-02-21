package responses

import (
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
	"github.com/recovery-flow/sso-oauth/resources"
)

func Account(account models.Account) resources.Account {
	return resources.Account{
		Data: resources.AccountData{
			Id:   account.ID.String(),
			Type: resources.AccountType,
			Attributes: resources.AccountDataAttributes{
				Email:     account.Email,
				Role:      string(account.Role),
				UpdatedAt: account.UpdatedAt,
				CreatedAt: account.CreatedAt,
			},
		},
	}
}
