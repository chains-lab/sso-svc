package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func (a App) TerminateByOwner(ctx context.Context, accountUD uuid.UUID) error {
	err := a.sessions.Terminate(ctx, accountUD)
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

func (a App) DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) error {
	if sessionID == initiatorSessionID {
		return fmt.Errorf("session can't be current")
	}
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

func (a App) TerminateByAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := a.accounts.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrAccountDoseNotExits
		default:
			return err
		}
	}

	if user.Role == roles.SuperUser {
		return ape.ErrSessionCannotDeleteForSuperUserByOtherUser
	}

	err = a.sessions.Terminate(ctx, userID)
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

func (a App) DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) error {
	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.ID == initiatorSessionID {
		return ape.ErrSessionCannotBeCurrent
	}

	if session.AccountID == initiatorID {
		return ape.ErrSessionCannotBeCurrentAccount
	}

	user, err := a.accounts.GetByID(ctx, session.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrAccountDoseNotExits
		default:
			return err
		}
	}

	if user.Role == roles.SuperUser {
		return ape.ErrSessionCannotDeleteForSuperUserByOtherUser
	}

	err = a.sessions.Delete(ctx, sessionID)
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

func (a App) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, ape.ErrSessionDoseNotExits
		default:
			return models.Session{}, err
		}
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
	_, err := a.accounts.GetByID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.ErrAccountDoseNotExits
		default:
			return nil, err
		}
	}
	sessions, err := a.sessions.GetByAccountID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.ErrSessionsForAccountDoseNotExits
		default:
			return nil, err
		}
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
