package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const UserCreatedEvent = "user.created"

type UserCreatedPayload struct {
	User models.User `json:"user"`
}

func (p Service) PublishUserCreated(
	ctx context.Context,
	user models.User,
) error {
	return p.publish(
		ctx,
		contracts.UsersTopicV1,
		user.ID.String(),
		contracts.Envelope[UserCreatedPayload]{
			Event:     UserCreatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: UserCreatedPayload{
				User: user,
			},
		},
	)
}
