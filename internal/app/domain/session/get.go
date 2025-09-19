package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Session) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.query.New().FilterID(sessionID).Get(ctx)
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
