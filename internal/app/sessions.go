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

func (a App) TerminateSessionsByOwner(ctx context.Context, accountUD uuid.UUID) *ape.Error {
	_, appError := a.GetAccountByID(ctx, accountUD)
	if appError != nil {
		return appError
	}

	err := a.sessions.Terminate(ctx, accountUD)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionsForAccountNotExist(err)
		default:
			return ape.ErrorInternalServer(err)
		}
	}
	return nil
}

func (a App) DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) *ape.Error {
	if sessionID == initiatorSessionID {
		return ape.ErrorSessionCannotBeCurrent(fmt.Errorf("cannot delete session that is the initiator session"))
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

func (a App) TerminateSessionsByAdmin(ctx context.Context, userID uuid.UUID) *ape.Error {
	user, appErr := a.GetAccountByID(ctx, userID)
	if appErr != nil {
		return appErr
	}

	if user.Role == roles.SuperUser {
		return ape.ErrorSessionCannotDeleteSuperUserByOther(fmt.Errorf("cannot delete superuser sessions"))
	}

	err := a.sessions.Terminate(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorSessionsForAccountNotExist(err)
		default:
			return ape.ErrorInternalServer(err)
		}
	}

	return nil
}

func (a App) DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) *ape.Error {
	session, appErr := a.GetSession(ctx, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.ID == initiatorSessionID {
		return ape.ErrorSessionCannotBeCurrent(fmt.Errorf("cannot delete session that is the initiator session"))
	}

	if session.AccountID == initiatorID {
		return ape.ErrorSessionCannotBeCurrentAccount(fmt.Errorf("cannot delete session that is the one of initiator session"))
	}

	user, appErr := a.GetAccountByID(ctx, session.AccountID)
	if appErr != nil {
		return appErr
	}

	if user.Role == roles.SuperUser {
		return ape.ErrorSessionCannotDeleteSuperUserByOther(fmt.Errorf("cannot delete superuser sessions"))
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

func (a App) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, *ape.Error) {
	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Session{}, ape.ErrorSessionDoesNotExist(sessionID, err)
		default:
			return models.Session{}, ape.ErrorInternalServer(err)
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

func (a App) GetSessions(ctx context.Context, accountID uuid.UUID) ([]models.Session, *ape.Error) {
	_, appErr := a.GetAccountByID(ctx, accountID)
	if appErr != nil {
		return nil, appErr
	}

	sessions, err := a.sessions.GetByAccountID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ape.ErrorSessionsForAccountNotExist(err)
		default:
			return nil, ape.ErrorInternalServer(err)
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
