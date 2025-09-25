package session

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
	page uint,
	size uint,
) (models.SessionCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := s.db.Sessions().Page(limit, offset).FilterUserID(userID)

	total, err := query.Count(ctx)
	if err != nil {
		return models.SessionCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("counting rows, cause: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return models.SessionCollection{}, errx.ErrorInternal.Raise(
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

	return models.SessionCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}
