package tools

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/recovery-flow/tokens"
)

func GetSessionAndUserID(ctx context.Context) (sessionID uuid.UUID, userID uuid.UUID, err error) {
	sessionID, ok := ctx.Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, uuid.Nil, errors.New("sessions not authenticated")
	}

	userID, ok = ctx.Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, uuid.Nil, errors.New("user not authenticated")
	}

	return sessionID, userID, nil
}
