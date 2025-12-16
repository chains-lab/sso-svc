package writer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountUsernameChangeEvent = "account.username.change"

type AccountUsernameChangePayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}

func (s Service) WriteAccountUsernameChanged(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return s.outbox.CreateOutboxEvent(
		ctx,
		contracts.Message{
			Topic:     contracts.AccountsTopicV1,
			EventType: AccountUsernameChangeEvent,
			Key:       account.ID.String(),
			Payload: AccountUsernameChangePayload{
				Account: account,
				Email:   email,
			},
		},
	)
}
