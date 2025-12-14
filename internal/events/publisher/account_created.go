package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountCreatedEvent = "account.created"

type AccountCreatedPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email,omitempty"`
}

func (p Service) PublishAccountCreated(
	ctx context.Context,
	account entity.Account,
	email string,
) error {
	return p.publish(
		ctx,
		contracts.AccountsTopicV1,
		account.ID.String(),
		contracts.Envelope[AccountCreatedPayload]{
			Event:     AccountCreatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: AccountCreatedPayload{
				Account: account,
				Email:   email,
			},
		},
	)
}
