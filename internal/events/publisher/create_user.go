package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/sso-svc/internal/events"
	"github.com/google/uuid"
)

func (p Service) CreateUser(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	email string,
) error {
	type UserCreated struct {
		UserID   uuid.UUID `json:"user_id"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
	}

	return p.publish(
		ctx,
		events.Users,
		userID.String(),
		events.Envelope[UserCreated]{
			Event:     "user.created",
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: UserCreated{
				UserID:   userID,
				Username: username,
				Email:    email,
			},
		},
	)
}
