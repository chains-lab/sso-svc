package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Session) DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", sessionID, userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}
	return nil
}
