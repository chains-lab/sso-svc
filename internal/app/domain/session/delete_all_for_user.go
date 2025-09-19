package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/chains-lab/sso-svc/internal/errx"
)

func (s Session) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	err := s.query.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorUserNotFound.Raise(
				fmt.Errorf("no sessions found for user %s, cause: %w", userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete sessions for user %s, cause: %w", userID, err),
			)
		}
	}
	return nil
}
