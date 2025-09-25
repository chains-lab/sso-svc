package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Refresh(ctx context.Context, oldRefreshToken string) (models.TokensPair, error) {
	tokenData, err := s.jwt.ParseRefreshClaims(oldRefreshToken)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to decrypt refresh token claims, cause: %w", err),
		)
	}

	userID, err := uuid.Parse(tokenData.Subject)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to parse user id from token claims, cause: %w", err),
		)
	}

	session, err := s.db.Sessions().FilterID(tokenData.SessionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.TokensPair{}, errx.ErrorSessionNotFound.Raise(
				fmt.Errorf("session with id: %s not found for user %s, cause: %w", tokenData.SessionID, userID, err),
			)
		default:
			return models.TokensPair{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get session with id: %s for user %s, cause: %w", tokenData.SessionID, userID, err),
			)
		}
	}

	refresh, err := s.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err),
		)
	}

	if refresh != oldRefreshToken {
		return models.TokensPair{}, errx.ErrorSessionTokenMismatch.Raise(
			fmt.Errorf("refresh token does not match for session %s and user %s, cause: %w", session.ID, userID, err),
		)
	}

	refresh, err = s.jwt.GenerateRefresh(userID, tokenData.SessionID, tokenData.Role, tokenData.Verified)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for user %s, cause: %w", userID, err),
		)
	}

	refreshCrypto, err := s.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for user %s, cause: %w", userID, err),
		)
	}

	access, err := s.jwt.GenerateAccess(userID, tokenData.SessionID, tokenData.Role)
	if err != nil {
		return models.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for user %s, cause: %w", userID, err),
		)
	}

	err = s.db.Sessions().FilterID(tokenData.SessionID).Update(ctx, schemas.UpdateSessionInput{
		Token:    &refreshCrypto,
		LastUsed: time.Now().UTC(),
	})

	return models.TokensPair{
		SessionID: tokenData.SessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
