package callback

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/sso-svc/internal/events"
	"github.com/google/uuid"
)

type AuthSvc interface {
	UpdateCompany(
		ctx context.Context,
		userID uuid.UUID,
		companyID *uuid.UUID,
		role *string,
	) error

	UpdatePassword(
		ctx context.Context,
		userID uuid.UUID,
		oldPassword, newPassword string,
	) error
}

type Service struct {
	auth AuthSvc
}

func NewService(auth AuthSvc) *Service {
	return &Service{
		auth: auth,
	}
}

func decodeEnvelope[T any](b []byte) (events.Envelope[T], error) {
	var env events.Envelope[T]
	err := json.Unmarshal(b, &env)
	return env, err
}
