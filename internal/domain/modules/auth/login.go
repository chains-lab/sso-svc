package auth

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) LoginByEmail(ctx context.Context, email, password string) (entity.TokensPair, error) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	err = s.checkAccountPassword(ctx, account.ID, password)
	if err != nil {
		return entity.TokensPair{}, err
	}

	return s.createSession(ctx, account)
}

func (s Service) LoginByUsername(ctx context.Context, username, password string) (entity.TokensPair, error) {
	account, err := s.GetAccountByUsername(ctx, username)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	err = s.checkAccountPassword(ctx, account.ID, password)
	if err != nil {
		return entity.TokensPair{}, err
	}

	return s.createSession(ctx, account)
}

func (s Service) LoginByGoogle(ctx context.Context, email string) (entity.TokensPair, error) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.TokensPair{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.TokensPair{}, err
	}

	return s.createSession(ctx, account)
}

func (s Service) checkAccountPassword(
	ctx context.Context,
	accountID uuid.UUID,
	password string,
) error {
	passData, err := s.db.GetAccountPassword(ctx, accountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get account password, cause: %w", err),
		)
	}

	if err = passData.CheckPasswordMatch(password); err != nil {
		return err
	}

	return nil
}

func (s Service) createSession(
	ctx context.Context,
	account entity.Account,
) (entity.TokensPair, error) {
	sessionID := uuid.New()

	pair, err := s.createTokensPair(sessionID, account)
	if err != nil {
		return entity.TokensPair{}, err
	}

	refreshTokenCrypto, err := s.jwt.EncryptRefresh(pair.Refresh)
	if err != nil {
		return entity.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to encrypt refresh token for account %s, cause: %w", account.ID, err),
		)
	}

	_, err = s.db.CreateSession(ctx, sessionID, account.ID, refreshTokenCrypto)
	if err != nil {
		return entity.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to createSession session for account %s, cause: %w", account.ID, err),
		)
	}

	email, err := s.GetAccountEmail(ctx, account.ID)
	if err != nil {
		return entity.TokensPair{}, err
	}

	err = s.event.WriteAccountLogin(ctx, account, email.Email)
	if err != nil {
		return entity.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish account login event for account %s: %w", account.ID, err),
		)
	}

	return entity.TokensPair{
		SessionID: pair.SessionID,
		Refresh:   pair.Refresh,
		Access:    pair.Access,
	}, nil
}

func (s Service) createTokensPair(
	sessionID uuid.UUID,
	account entity.Account,
) (entity.TokensPair, error) {
	access, err := s.jwt.GenerateAccess(account, sessionID)
	if err != nil {
		return entity.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate access token for account %s, cause: %w", account.ID, err),
		)
	}

	refresh, err := s.jwt.GenerateRefresh(account, sessionID)
	if err != nil {
		return entity.TokensPair{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to generate refresh token for account %s, cause: %w", account.ID, err),
		)
	}

	return entity.TokensPair{
		SessionID: sessionID,
		Refresh:   refresh,
		Access:    access,
	}, nil
}
