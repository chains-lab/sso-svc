package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/pagination"
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
			return models.Session{}, models.TokensPair{}, errx.RaiseSessionNotFound(
				ctx,
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
				sessionID,
				userID,
			)
		default:
			return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
		}
	}

	if session.Client != client {
		return models.Session{}, models.TokensPair{}, errx.RaiseSessionClientMismatch(
			ctx,
			fmt.Errorf("client mismatch"),
		)
	}

	access, err := a.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	oldRefresh, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	if oldRefresh != token {
		return models.Session{}, models.TokensPair{}, errx.RaiseSessionTokenMismatch(
			ctx,
			fmt.Errorf("refresh token mismatch"),
		)
	}

	newRefresh, err := a.jwt.GenerateRefresh(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	LastUsed := time.Now().UTC()

	err = a.sessionQ.New().FilterID(sessionID).Update(ctx, map[string]any{
		"token":     refreshCrypto,
		"ip":        ip,
		"last_used": LastUsed,
	})
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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
			return models.Session{}, errx.RaiseSessionNotFound(
				ctx,
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
				sessionID,
				userID,
			)
		default:
			return models.Session{}, errx.RaiseInternal(ctx, err)
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

func (a App) GetUserSessions(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error) {
	limit, offset := pagination.CalculateLimitOffset(pag)

	sessions, err := a.sessionQ.New().FilterID(userID).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, errx.RaiseSessionsForUserNotFound(
				ctx, fmt.Errorf("no sessions found for user %s", userID),
			)
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	total, err := a.sessionQ.New().FilterUserID(userID).Count(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, errx.RaiseSessionsForUserNotFound(
				ctx, fmt.Errorf("no sessions found for user %s", userID),
			)
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	result := make([]models.Session, len(sessions))
	for i, session := range sessions {
		result[i] = models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			IP:        session.IP,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return result, pagination.Response{
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
			return errx.RaiseSessionNotFound(
				ctx,
				fmt.Errorf("session with id: %s not found", sessionID),
				sessionID,
				userID,
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}
	return nil
}

func (a App) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := a.sessionQ.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("no sessions found for user %s", userID),
				userID,
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}
	return nil
}
