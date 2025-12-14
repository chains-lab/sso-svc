package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountPasswordChangeEvent = "account.password.change"

type AccountPasswordChangePayload struct {
	Account entity.Account `json:"account"`
}

func (p Service) PublishAccountPasswordChanged(
	ctx context.Context,
	account entity.Account,
) error {
	return p.publish(
		ctx,
		contracts.AccountsTopicV1,
		account.ID.String(),
		contracts.Envelope[AccountPasswordChangePayload]{
			Event:     AccountPasswordChangeEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: AccountPasswordChangePayload{
				Account: account,
			},
		},
	)
}
