package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountLoginEvent = "account.login"

type AccountLoginPayload struct {
	Account entity.Account `json:"omitempty"`
}

func (p Service) PublishAccountLogin(
	ctx context.Context,
	account entity.Account,
) error {
	return p.publish(
		ctx,
		contracts.AccountsTopicV1,
		account.ID.String(),
		contracts.Envelope[AccountLoginPayload]{
			Event:     AccountLoginEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: AccountLoginPayload{
				Account: account,
			},
		},
	)
}
