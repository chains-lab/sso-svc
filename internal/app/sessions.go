package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, ip, token string) (models.Session, models.TokensPair, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	session, err := a.sessionQ.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, models.TokensPair{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
			)
		default:
			return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}

	if session.Client != client {
		return models.Session{}, models.TokensPair{}, errx.ErrorSessionClientMismatch.Raise(
			fmt.Errorf("client mismatch"),
		)
	}

	access, err := a.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err),
		)
	}

	oldRefresh, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to decrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	if oldRefresh != token {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("refresh token mismatch"),
		)
	}

	newRefresh, err := a.jwt.GenerateRefresh(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err),
		)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	LastUsed := time.Now().UTC()

	err = a.sessionQ.New().FilterID(sessionID).Update(ctx, map[string]any{
		"token":     refreshCrypto,
		"ip":        ip,
		"last_used": LastUsed,
	})
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update session for user %s: %w", user.ID, err),
		)
	}

	return models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			IP:        ip,
			LastUsed:  LastUsed,
			CreatedAt: session.CreatedAt,
		}, models.TokensPair{
			Refresh: newRefresh,
			Access:  access,
		}, nil
}

func (a App) GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, err := a.sessionQ.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
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

func (a App) GetUserSessions(
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

	query := a.sessionQ.New().Page(limit, offset).FilterUserID(userID)

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

func (a App) DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := a.sessionQ.New().FilterID(sessionID).Delete(ctx)
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

func (a App) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := a.sessionQ.New().FilterUserID(userID).Delete(ctx)
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
