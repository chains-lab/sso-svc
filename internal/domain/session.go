package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) CreateTokensPair(
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

func (s Service) CreateSession(
	ctx context.Context,
	account entity.Account,
) (entity.TokensPair, error) {
	sessionID := uuid.New()

	pair, err := s.CreateTokensPair(sessionID, account)
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
			fmt.Errorf("failed to CreateSession session for account %s, cause: %w", account.ID, err),
		)
	}

	email, err := s.GetAccountEmailData(ctx, account.ID)
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

func (s Service) GetSessionForAccount(ctx context.Context, accountID, sessionID uuid.UUID) (entity.Session, error) {
	account, err := s.GetAccountByID(ctx, accountID)
	if err != nil {
		return entity.Session{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.Session{}, err
	}

	session, err := s.db.GetAccountSession(ctx, accountID, sessionID)
	if err != nil {
		return entity.Session{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get session with id: %s for account %s, cause: %w", sessionID, accountID, err),
		)
	}

	if session.IsNil() {
		return entity.Session{}, errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session with id: %s for account %s not found", sessionID, accountID),
		)
	}

	return session, nil
}

func (s Service) GetSessionsForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	page int32,
	size int32,
) (entity.SessionsCollection, error) {
	account, err := s.GetAccountByID(ctx, accountID)
	if err != nil {
		return entity.SessionsCollection{}, err
	}

	if err = account.CanInteract(); err != nil {
		return entity.SessionsCollection{}, err
	}

	sessions, err := s.db.GetSessionsForAccount(ctx, accountID, page, size)
	if err != nil {
		return entity.SessionsCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list sessions for account %s, cause: %w", accountID, err),
		)
	}

	return sessions, nil
}

func (s Service) Logout(ctx context.Context, accountID, sessionID uuid.UUID) error {
	err := s.db.DeleteAccountSession(ctx, accountID, sessionID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete session with id: %s, cause: %w", sessionID, err),
		)
	}

	return nil
}

func (s Service) DeleteOwnSession(ctx context.Context, accountID, sessionID uuid.UUID) error {
	account, err := s.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	if err = account.CanInteract(); err != nil {
		return err
	}

	err = s.db.DeleteAccountSession(ctx, accountID, sessionID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete session with id: %s for account %s, cause: %w", sessionID, accountID, err),
		)
	}

	return nil
}

func (s Service) DeleteOwnSessions(ctx context.Context, accountID uuid.UUID) error {
	account, err := s.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	if err = account.CanInteract(); err != nil {
		return err
	}

	err = s.db.DeleteSessionsForAccount(ctx, accountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete sessions for account %s, cause: %w", accountID, err),
		)
	}

	return nil
}
