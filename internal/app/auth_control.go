package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
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
			err = a.AccountCreate(ctx, request.Email, roles.User)
			if err != nil {
				return err
			}

			account, err = a.accounts.GetByEmail(ctx, request.Email)
			if err != nil {
				return err
			}
		}

		refresh, err := a.jwt.GenerateRefresh(account.ID, sessionID, account.ID, account.Role)
		if err != nil {
			return err
		}

		access, err := a.jwt.GenerateAccess(account.ID, sessionID, account.ID, account.Role)
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
			CreatedAt: time.Now().UTC(),
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

func (a App) Refresh(ctx context.Context, accountID, sessionID uuid.UUID, request RefreshRequest) (models.Session, error) {
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
		return models.Session{}, ape.ErrSessionsClientMismatch
	}

	refreshToken, err := a.jwt.DecryptRefresh(session.Token)
	if err != nil {
		return models.Session{}, err
	}

	if refreshToken != request.Token {
		return models.Session{}, ape.ErrSessionsTokenMismatch
	}

	access, err := a.jwt.GenerateAccess(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, err
	}

	refresh, err := a.jwt.GenerateRefresh(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, err
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, err
	}

	err = a.sessions.Update(ctx, sessionID, repo.SessionUpdateRequest{
		Token:    &refreshCrypto,
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

	err := a.sessions.Delete(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrSessionDoseNotExits
		default:
			return err
		}
	}
	return nil
}
