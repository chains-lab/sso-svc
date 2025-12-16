package writer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountLoginEvent = "account.login"

type AccountLoginPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}

func (s Service) WriteAccountLogin(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return s.repo.CreateOutboxEvent(
		ctx,
		contracts.Event{
			Topic:     contracts.AccountsTopicV1,
			EventType: AccountLoginEvent,
			Key:       account.ID.String(),
			Payload: AccountLoginPayload{
				Account: account,
				Email:   email,
			},
		},
	)
}
