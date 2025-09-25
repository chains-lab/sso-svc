package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(ctx context.Context, email, password string) (models.TokensPair, error) {
	emailData, err := s.db.UsersEmail().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	user, err := s.db.Users().FilterID(emailData.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	if user.Status == enum.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	secret, err := s.db.UsersPassword().FilterID(user.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("password for user %s not found, cause: %w", user.ID, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting password for user %s, cause: %w", user.ID, err),
			)
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.TokensPair{}, errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials for user %s, cause: %w", user.ID, err),
			)
		}

		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash for user %s, cause: %w", user.ID, err),
		)
	}

	pair, err := s.Create(ctx, user.ID, user.Role, emailData.Verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to Create session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (s Service) LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error) {
	emailData, err := s.db.UsersEmail().FilterEmail(email).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	user, err := s.db.Users().FilterID(emailData.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
				fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
			)
		}
	}

	pair, err := s.Create(ctx, user.ID, user.Role, emailData.Verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to Create session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (s Service) Create(
	ctx context.Context,
	userID uuid.UUID,
	role string,
	verified bool,
) (models.TokensPair, error) {
	sessionID := uuid.New()

	access, err := s.jwt.GenerateAccess(userID, sessionID, role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", userID, err))
	}

	refresh, err := s.jwt.GenerateRefresh(userID, sessionID, role, verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err))
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", userID, err))
	}

	session := schemas.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshCrypto,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = s.db.Sessions().Insert(ctx, session)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to Create session for user %s, cause: %w", userID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
