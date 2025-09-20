package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Session) DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).FilterID(sessionID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete session with id: %s for user %s, cause: %w", sessionID, userID, err),
		)
	}
	return nil
}
