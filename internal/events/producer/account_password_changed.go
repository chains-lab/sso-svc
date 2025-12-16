package producer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountPasswordChangeEvent = "account.password.change"

type AccountPasswordChangePayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}

func (s Service) WriteAccountPasswordChanged(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return s.outbox.CreateOutboxEvent(
		ctx,
		contracts.Message{
			Topic:     contracts.AccountsTopicV1,
			EventType: AccountPasswordChangeEvent,
			Key:       account.ID.String(),
			Payload: AccountPasswordChangePayload{
				Account: account,
				Email:   email,
			},
		},
	)
}
