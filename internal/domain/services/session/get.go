package session

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.db.GetSession(ctx, sessionID)
	if err != nil {
		return models.Session{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get session with id: %s cause: %w", sessionID, err),
		)
	}

	if session == (schemas.Session{}) {
		return models.Session{}, errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session with id: %s not found", sessionID),
		)
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Service) GetForUser(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := s.db.GetOneSessionForUser(ctx, userID, sessionID)
	if err != nil {
		return models.Session{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", sessionID, userID, err),
		)
	}
	if session == (schemas.Session{}) {
		return models.Session{}, errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session with id: %s for user %s not found", sessionID, userID),
		)
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (s Service) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
	page uint,
	size uint,
) (models.SessionsCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	rows, total, err := s.db.GetAllSessionsForUser(ctx, userID, limit, offset)
	if err != nil {
		return models.SessionsCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows, cause: %w", err),
		)
	}

	result := make([]models.Session, len(rows))
	for i, session := range rows {
		result[i] = models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return models.SessionsCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}
