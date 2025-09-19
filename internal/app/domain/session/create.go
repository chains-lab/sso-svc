package session

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Session) Create(
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

	session := dbx.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshCrypto,
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	err = s.query.Insert(ctx, session)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create session for user %s, cause: %w", userID, err),
		)
	}

	return models.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
