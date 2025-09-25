package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.db.Sessions().
		FilterID(sessionID).Delete(ctx)
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

func (s Service) DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.db.Sessions().FilterUserID(userID).FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s for user %s not found, cause: %w", sessionID, userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}

	return nil
}

func (s Service) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	err := s.db.Sessions().FilterUserID(userID).Delete(ctx)
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
