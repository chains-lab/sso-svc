package responses

import (
	"github.com/hs-zavet/sso-oauth/resources"
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
