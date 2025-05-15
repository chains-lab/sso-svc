package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (a App) Login(ctx context.Context, request LoginRequest) (models.Session, *ape.Error) {
	sessionID := uuid.New()
	var result models.Session

	txErr := a.accounts.Transaction(func(ctx context.Context) error {

		account, err := a.accounts.GetByEmail(ctx, request.Email)
		if err != nil {
			ID := uuid.New()
			CreatedAt := time.Now().UTC()

			err := a.accounts.Create(ctx, repo.AccountCreateRequest{
				ID:           ID,
				Email:        request.Email,
				Role:         roles.User,
				Subscription: uuid.Nil,
				CreatedAt:    CreatedAt,
			})
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

	if txErr != nil {
		switch {
		default:
			return models.Session{}, ape.ErrorInternalServer(txErr)
		}
	}

	return result, nil
}

type RefreshRequest struct {
	Token  string `json:"token"`
	Client string `json:"client"`
}

func (a App) Refresh(ctx context.Context, accountID, sessionID uuid.UUID, request RefreshRequest) (models.Session, *ape.Error) {
	LastUsed := time.Now().UTC()

	session, appErr := a.GetSession(ctx, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	account, appErr := a.GetAccountByID(ctx, session.AccountID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	if session.Client != request.Client {
		return models.Session{}, ape.ErrorSessionClientMismatch(fmt.Errorf("session client mismatch"))
	}

	refreshToken, err := a.jwt.DecryptRefresh(session.Refresh)
	if err != nil {
		return models.Session{}, ape.ErrorInternalServer(err)
	}

	if refreshToken != request.Token {
		return models.Session{}, ape.ErrorSessionTokenMismatch(err)
	}

	access, err := a.jwt.GenerateAccess(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, ape.ErrorInternalServer(err)
	}

	refresh, err := a.jwt.GenerateRefresh(session.AccountID, session.ID, account.Subscription, account.Role)
	if err != nil {
		return models.Session{}, ape.ErrorInternalServer(err)
	}

	refreshCrypto, err := a.jwt.EncryptRefresh(refresh)
	if err != nil {
		return models.Session{}, ape.ErrorInternalServer(err)
	}

	err = a.sessions.Update(ctx, sessionID, repo.SessionUpdateRequest{
		Token:    &refreshCrypto,
		LastUsed: LastUsed,
	})
	if err != nil {
		return models.Session{}, ape.ErrorInternalServer(err)
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

func (a App) Logout(ctx context.Context, sessionID uuid.UUID) *ape.Error {
	_, appErr := a.GetSession(ctx, sessionID)
	if appErr != nil {
		return appErr
	}

	err := a.sessions.Delete(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionDoesNotExist(sessionID, err)
		default:
			return ape.ErrorInternalServer(err)
		}
	}
	return nil
}
