package writer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountPasswordChangeEvent = "account.password.change"

type AccountPasswordChangePayload struct {
	Account entity.Account `json:"account"`
}

func (s Service) WriteAccountPasswordChanged(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return s.addToOutbox(
		ctx,
		contracts.Event{
			Topic:     contracts.AccountsTopicV1,
			EventType: AccountPasswordChangeEvent,
			Key:       account.ID.String(),
			Payload: AccountPasswordChangePayload{
				Account: account,
			},
		},
	)
}
