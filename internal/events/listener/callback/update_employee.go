package callback

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type EmployeeUpdate struct {
	UserID    uuid.UUID  `json:"user_id"`
	CompanyID *uuid.UUID `json:"company_id"`
	Role      *string    `json:"role"`
}

func (s Service) UpdateEmployee(ctx context.Context, event kafka.Message) error {
	env, err := decodeEnvelope[EmployeeUpdate](event.Value)
	if err != nil {
		return err
	}

	if err = s.auth.UpdateCompany(ctx, env.Data.UserID, env.Data.CompanyID, env.Data.Role); err != nil {
		return err
	}

	return nil
}
