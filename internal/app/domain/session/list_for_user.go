package session

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Session) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	query := s.query.New().Page(limit, offset).FilterUserID(userID)

	for _, sort := range sort {
		switch sort.Field {
		case "created_at":
			query = query.OrderCreatedAt(sort.Ascend)
		default:

		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("counting rows, cause: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
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

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}
