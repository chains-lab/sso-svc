package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/sso-oauth/internal/repo"
	"github.com/hs-zavet/tokens/identity"
)

type LoginRequest struct {
	Email  string
	Client string
}

func (a App) Login(ctx context.Context, request LoginRequest) (models.Session, error) {
	sessionID := uuid.New()
	var result models.Session

	err := a.accounts.Transaction(func(ctx context.Context) error {
		account, err := a.accounts.GetByEmail(ctx, request.Email)
		if err != nil {
			err = a.AccountCreate(ctx, request.Email)
			if err != nil {
				return err
			}

			account, err = a.accounts.GetByEmail(ctx, request.Email)
			if err != nil {
				return err
			}
		}

		role, err := identity.ParseIdentityType(account.Role)
		if err != nil {
			return err
		}

		refresh, err := a.jwt.GenerateRefresh(account.ID, sessionID, account.ID, role)
		if err != nil {
			return err
		}

		access, err := a.jwt.GenerateAccess(account.ID, sessionID, account.ID, role)
		if err != nil {
			return err
		}

		refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
		if err != nil {
			return err
		}

		err = a.sessions.Create(ctx, repo.SessionCreateRequest{
			ID:        sessionID,
			AccountID: account.ID,
			Token:     refreshCrypto,
			Client:    request.Client,
		})
		if err != nil {
			return err
		}

		now := time.Now().UTC()

		result = models.Session{
			ID:        sessionID,
			AccountID: account.ID,
			Access:    access,
			Refresh:   refresh,
			Client:    request.Client,
			LastUsed:  now,
			CreatedAt: now,
		}
		return nil
	})

	if err != nil {
		return models.Session{}, err
	}

	return result, nil
}

type RefreshRequest struct {
	Token  string `json:"token"`
	Client string `json:"client"`
}

func (a App) Refresh(ctx context.Context, sessionID uuid.UUID, request RefreshRequest) (models.Session, error) {
	LastUsed := time.Now().UTC()

	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return models.Session{}, err
	}

	account, err := a.accounts.GetByID(ctx, session.AccountID)
	if err != nil {
		return models.Session{}, err
	}

	if session.Client != request.Client {
		return models.Session{}, fmt.Errorf("client does not match")
	}

	refreshToken, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, err
	}

	if refreshToken != request.Token {
		return models.Session{}, fmt.Errorf("token does not match")
	}

	role, err := identity.ParseIdentityType(account.Role)
	if err != nil {
		return models.Session{}, err
	}

	access, err := a.jwt.GenerateAccess(session.AccountID, session.ID, account.Subscription, role)
	if err != nil {
		return models.Session{}, err
	}

	refresh, err := a.jwt.GenerateRefresh(session.AccountID, session.ID, account.Subscription, role)
	if err != nil {
		return models.Session{}, err
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, err
	}

	err = a.sessions.Update(ctx, sessionID, repo.SessionUpdateRequest{
		Token:    refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, err
	}

	return models.Session{
		ID:        session.ID,
		AccountID: session.AccountID,
		Access:    access,
		Refresh:   refresh,
		Client:    session.Client,
		LastUsed:  LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (a App) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return a.sessions.Delete(ctx, sessionID)
}

func (a App) Terminate(ctx context.Context, sessionID uuid.UUID) error {
	return a.sessions.Terminate(ctx, sessionID)
}

func (a App) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return a.sessions.Delete(ctx, sessionID)
}

func (a App) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return models.Session{}, err
	}

	return models.Session{
		ID:        session.ID,
		AccountID: session.AccountID,
		Client:    session.Client,
		LastUsed:  session.LastUsed,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (a App) GetSessions(ctx context.Context, accountID uuid.UUID) ([]models.Session, error) {
	sessions, err := a.sessions.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Session, len(sessions))
	for i, session := range sessions {
		result[i] = models.Session{
			ID:        session.ID,
			AccountID: session.AccountID,
			Client:    session.Client,
			LastUsed:  session.LastUsed,
			CreatedAt: session.CreatedAt,
		}
	}

	return result, nil
}
