package session

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := s.db.DeleteSession(ctx, sessionID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete session with id: %s, cause: %w", sessionID, err),
		)
	}

	return nil
}

func (s Service) DeleteOneForUser(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := s.db.DeleteOneSessionForUser(ctx, userID, sessionID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete session with id: %s for user %s, cause: %w", sessionID, userID, err),
		)
	}

	return nil
}

func (s Service) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	err := s.db.DeleteAllSessionsForUser(ctx, userID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete sessions for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
