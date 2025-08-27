package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a App) GoogleLogin(ctx context.Context, email, client, ip string) (models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err))
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err))
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", user.ID, err))
	}

	session := dbx.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     refreshCrypto,
		Client:    client,
		IP:        ip,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = a.sessionQ.New().Insert(ctx, session)
	if err != nil {
		switch {
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
			)
		}
	}

	return models.TokensPair{
		Refresh: refresh,
		Access:  access,
	}, nil
}

func (a App) Login(ctx context.Context, email, password, client, ip string) (models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, err
	}

	if user.Status == constant.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	secret, err := a.passQ.New().FilterID(user.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("password for user %s not found: %w", user.ID, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting password for user %s: %w", user.ID, err),
			)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("invalid credentials for user %s: %w", user.ID, err),
			)
		}
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash for user %s: %w", user.ID, err),
		)
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err),
		)
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err),
		)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	session := dbx.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     refreshCrypto,
		Client:    client,
		IP:        ip,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = a.sessionQ.New().Insert(ctx, session)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("creating session for user %s: %w", user.ID, err),
			)
		}
	}

	return models.TokensPair{
		Refresh: refresh,
		Access:  access,
	}, nil
}

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, ip, token string) (models.TokensPair, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.TokensPair{}, appErr
	}

	if user.Status == constant.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	session, err := a.sessionQ.New().FilterID(sessionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s", sessionID, userID),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s: %w", sessionID, userID, err),
			)
		}
	}

	if session.Client != client {
		return models.TokensPair{}, errx.ErrorSessionClientMismatch.Raise(
			fmt.Errorf("client mismatch"),
		)
	}

	access, err := a.jwt.GenerateAccess(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s: %w", user.ID, err),
		)
	}

	oldRefresh, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to decrypt refresh token for user %s: %w", user.ID, err),
		)
	}

	if oldRefresh != token {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("refresh token mismatch"),
		)
	}

	newRefresh, err := a.jwt.GenerateRefresh(session.UserID, session.ID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s: %w", user.ID, err),
		)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(newRefresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
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
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   newRefresh,
		Access:    access,
	}, nil
}
