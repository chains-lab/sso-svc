package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const UserPasswordChangeEvent = "user.password.change"

type UserPasswordChangePayload struct {
	User models.User `json:"user"`
}

func (p Service) PublishUserPasswordChanged(
	ctx context.Context,
	user models.User,
) error {
	return p.publish(
		ctx,
		contracts.UsersTopicV1,
		user.ID.String(),
		contracts.Envelope[UserPasswordChangePayload]{
			Event:     UserPasswordChangeEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: UserPasswordChangePayload{
				User: user,
			},
		},
	)
}
