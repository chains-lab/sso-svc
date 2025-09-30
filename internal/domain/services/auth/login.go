package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/infra/password"
	"github.com/google/uuid"
)

func (s Service) Login(ctx context.Context, email, pass string) (models.TokensPair, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
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

	if err = password.CheckPasswordMatch(pass, user.PasswordHash); err != nil {
		return models.TokensPair{}, err
	}

	pair, err := s.CreateSession(ctx, user.ID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to CreateSession session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (s Service) LoginByGoogle(ctx context.Context, email string) (models.TokensPair, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
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

	pair, err := s.CreateSession(ctx, user.ID, user.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to CreateSession session for user %s: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (s Service) CreateSession(
	ctx context.Context,
	userID uuid.UUID,
	role string,
) (models.TokensPair, error) {
	sessionID := uuid.New()

	access, err := s.jwt.GenerateAccess(userID, sessionID, role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", userID, err))
	}

	refresh, err := s.jwt.GenerateRefresh(userID, sessionID, role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err))
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", userID, err))
	}

	err = s.db.CreateSession(ctx, schemas.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshCrypto,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to CreateSession session for user %s, cause: %w", userID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
