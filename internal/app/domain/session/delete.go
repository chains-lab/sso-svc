package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Session) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.query.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found, cause: %w", sessionID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s, cause: %w", sessionID, err),
			)
		}
	}
	return nil
}
