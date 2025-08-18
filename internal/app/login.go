package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a App) GoogleLogin(ctx context.Context, email, client, ip string) (models.Session, models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.Session{}, models.TokensPair{}, err
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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
			return models.Session{}, models.TokensPair{}, errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
				user.ID,
			)
		default:
			return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			IP:        session.IP,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}, models.TokensPair{
			Refresh: refresh,
			Access:  access,
		}, nil
}

func (a App) Login(ctx context.Context, email, password, client, ip string) (models.Session, models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.Session{}, models.TokensPair{}, err
	}

	secret, err := a.passQ.New().FilterID(user.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, models.TokensPair{}, errx.RaiseUserNotFoundByEmail(ctx, err, email)
		default:
			return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.Session{}, models.TokensPair{}, errx.RaiseLoginIsIncorrect(
				ctx,
				fmt.Errorf("invalid credentials for user %s: %w", user.ID, err),
			)
		}
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	sessionID := uuid.New()

	access, err := a.jwt.GenerateAccess(user.ID, sessionID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refresh, err := a.jwt.GenerateRefresh(user.ID, sessionID, user.Role)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
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
			return models.Session{}, models.TokensPair{}, errx.RaiseUserNotFound(
				ctx,
				fmt.Errorf("failed to create session for user %s: %w", user.ID, err),
				user.ID,
			)
		default:
			return models.Session{}, models.TokensPair{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			Client:    session.Client,
			IP:        session.IP,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}, models.TokensPair{
			Refresh: refresh,
			Access:  access,
		}, nil
}
