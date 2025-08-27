package meta

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserData struct {
	ID        uuid.UUID `json:"sub,omitempty"`
	SessionID uuid.UUID `json:"session_id,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func User(ctx context.Context) (UserData, error) {
	if ctx == nil {
		return UserData{}, fmt.Errorf("mising context")
	}

	userData, ok := ctx.Value(UserCtxKey).(UserData)
	if !ok {
		return UserData{}, fmt.Errorf("mising context")
	}

	return userData, nil
}
