package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/domain/enum"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Login(ctx context.Context, email, pass string) (models.TokensPair, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
		)
	}

	if user.IsNil() {
		return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
			fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
		)
	}

	if user.Status == enum.UserStatusBlocked {
		return models.TokensPair{}, errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("initator user %s is blocked", user.ID),
		)
	}

	passData, err := s.db.GetUserPassword(ctx, user.ID)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user password, cause: %w", err),
		)
	}

	if err = s.pass.CheckPasswordMatch(pass, passData.Hash); err != nil {
		return models.TokensPair{}, err
	}

	pair, err := s.CreateSession(ctx, user)
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
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user with email '%s', cause: %w", email, err),
		)
	}

	if user.IsNil() {
		return models.TokensPair{}, errx.ErrorUserNotFound.Raise(
			fmt.Errorf("user with email '%s' not found, cause: %w", email, err),
		)
	}

	pair, err := s.CreateSession(ctx, user)
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
	user models.User,
) (models.TokensPair, error) {
	sessionID := uuid.New()

	access, err := s.jwt.GenerateAccess(user, sessionID)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", user.ID, err))
	}

	refresh, err := s.jwt.GenerateRefresh(user, sessionID)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", user.ID, err))
	}

	refreshTokenCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", user.ID, err),
		)
	}

	err = s.db.CreateSession(ctx, models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}, refreshTokenCrypto)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to CreateSession session for user %s, cause: %w", user.ID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
