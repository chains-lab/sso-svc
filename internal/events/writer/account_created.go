package writer

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountCreatedEvent = "account.created"

type AccountCreatedPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email,omitempty"`
}

func (s Service) WriteAccountCreated(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return s.repo.CreateOutboxEvent(
		ctx,
		contracts.Event{
			Topic:     contracts.AccountsTopicV1,
			EventType: AccountCreatedEvent,
			Key:       account.ID.String(),
			Payload: AccountCreatedPayload{
				Account: account,
				Email:   email,
			},
		},
	)
}
