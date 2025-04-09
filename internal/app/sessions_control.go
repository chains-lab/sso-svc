package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/tokens/roles"
)

func (a App) TerminateByOwner(ctx context.Context, accountUD uuid.UUID) error {
	return a.sessions.Terminate(ctx, accountUD)
}

func (a App) DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) error {
	if sessionID == initiatorSessionID {
		return fmt.Errorf("session can't be current")
	}
	return a.sessions.Delete(ctx, sessionID)
}

func (a App) TerminateByAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := a.accounts.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role == roles.SuperUser {
		return fmt.Errorf("cannot delete superuser")
	}

	return a.sessions.Terminate(ctx, userID)
}

func (a App) DeleteSessionByAdmin(ctx context.Context, sessionID, initiatorID, initiatorSessionID uuid.UUID) error {
	session, err := a.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.ID == initiatorSessionID {
		return fmt.Errorf("session can't be current")
	}

	if session.AccountID == initiatorID {
		return fmt.Errorf("account can't be current")
	}

	user, err := a.accounts.GetByID(ctx, session.AccountID)
	if err != nil {
		return err
	}

	if user.Role == roles.SuperUser {
		return fmt.Errorf("cannot delete superuser")
	}

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
