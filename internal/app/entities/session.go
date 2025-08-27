package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

type Session struct {
	query dbx.SessionsQ
}

func (a Session) GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := a.query.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
			)
		default:
			return models.Session{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Client:    session.Client,
		IP:        session.IP,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (a Session) GetUserSessions(
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

	query := a.query.New().Page(limit, offset).FilterUserID(userID)

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "created_at":
			query = query.OrderCreatedAt(ascend)
		default:

		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("counting rows: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows: %w", err),
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
			Client:    session.Client,
			IP:        session.IP,
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

func (a Session) DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := a.query.New().FilterID(sessionID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found", sessionID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}
	return nil
}

func (a Session) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := a.query.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("no sessions found for user %s", userID),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete sessions for user %s: %w", userID, err),
			)
		}
	}
	return nil
}
