package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

const UserLoginEvent = "user.login"

type UserLoginPayload struct {
	User models.User `json:"user"`
}

func (p Service) PublishUserLogin(
	ctx context.Context,
	user models.User,
) error {
	return p.publish(
		ctx,
		contracts.UsersTopicV1,
		user.ID.String(),
		contracts.Envelope[UserLoginPayload]{
			Event:     UserLoginEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: UserLoginPayload{
				User: user,
			},
		},
	)
}
