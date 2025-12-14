package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const AccountUsernameChangeEvent = "account.username.change"

type AccountUsernameChangePayload struct {
	Account entity.Account `json:"account"`
}

func (p Service) PublishAccountUsernameChanged(
	ctx context.Context,
	account entity.Account,
) error {
	return p.publish(
		ctx,
		contracts.AccountsTopicV1,
		account.ID.String(),
		contracts.Envelope[AccountUsernameChangePayload]{
			Event:     AccountUsernameChangeEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: AccountUsernameChangePayload{
				Account: account,
			},
		},
	)
}
