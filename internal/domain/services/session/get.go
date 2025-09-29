package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/models"

	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.db.Sessions().FilterID(sessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found, cause: %w", sessionID, err),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s cause: %w", sessionID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Service) GetForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.db.Sessions().FilterUserID(userID).FilterID(sessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", sessionID, userID, err),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", sessionID, userID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}
