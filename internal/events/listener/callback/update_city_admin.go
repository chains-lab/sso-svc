package callback

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type CityAdminUpdate struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

func (s Service) UpdateCityAdmin(ctx context.Context, event kafka.Message) error {
	env, err := decodeEnvelope[CityAdminUpdate](event.Value)
	if err != nil {
		return err
	}

	role := env.Data.Role
	if err := s.auth.UpdateCompany(ctx, env.Data.UserID, nil, &role); err != nil {
		return err
	}

	return nil
}
