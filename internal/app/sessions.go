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

func (a App) TerminateSessionsByOwner(ctx context.Context, userID uuid.UUID) *ape.Error {
	_, appError := a.GetUserByID(ctx, userID)
	if appError != nil {
		return appError
	}

	appErr := a.sessions.Terminate(ctx, userID)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a App) DeleteSessionByOwner(ctx context.Context, sessionID, initiatorSessionID uuid.UUID) *ape.Error {
	if sessionID == initiatorSessionID {
		return ape.ErrorSessionCannotBeCurrent(fmt.Errorf("cannot delete session that is the initiator session"))
	}

	appErr := a.sessions.Delete(ctx, sessionID)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a App) TerminateSessionsByAdmin(ctx context.Context, userID uuid.UUID) *ape.Error {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return appErr
	}

	if user.Role == roles.SuperUser {
		return ape.ErrorSessionCannotDeleteSuperUserByOther(fmt.Errorf("cannot delete superuser sessions"))
	}

	appErr = a.sessions.Terminate(ctx, userID)
	if appErr != nil {
		return appErr
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

	if session.UserID == initiatorID {
		return ape.ErrorSessionCannotBeCurrentUser(fmt.Errorf("cannot delete session that is the one of initiator session"))
	}

	user, appErr := a.GetUserByID(ctx, session.UserID)
	if appErr != nil {
		return appErr
	}

	if user.Role == roles.SuperUser {
		return ape.ErrorSessionCannotDeleteSuperUserByOther(fmt.Errorf("cannot delete superuser sessions"))
	}

	appErr = a.sessions.Delete(ctx, sessionID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, *ape.Error) {
	session, appErr := a.sessions.Get(ctx, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, *ape.Error) {
	_, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	sessions, appErr := a.sessions.GetByUserID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	return sessions, nil
}

func (a App) Login(ctx context.Context, email, client string) (models.Session, *ape.Error) {
	user, appErr := a.users.GetByEmail(ctx, email)
	if appErr != nil {

		//Registration flow
		if errors.Is(appErr.Unwrap(), sql.ErrNoRows) {
			appErr = a.users.Create(ctx, email, roles.User)
			if appErr != nil {
				switch {
				case errors.Is(appErr.Unwrap(), sql.ErrNoRows):
					return models.Session{}, ape.ErrorUserDoesNotExistByEmail(email, appErr.Unwrap())
				default:
					return models.Session{}, ape.ErrorInternal(appErr.Unwrap())
				}
			}

			user, appErr = a.users.GetByEmail(ctx, email)
			if appErr != nil {

				// It a good return internal error here anyway, because we already created the user in logic above
				return models.Session{}, ape.ErrorInternal(appErr.Unwrap())
			}

			session, appErr := a.sessions.Create(ctx, user, client)
			if appErr != nil {

				// If we fail to create a session after creating an user, we should return an internal error
				return models.Session{}, ape.ErrorInternal(appErr.Unwrap())
			}

			return session, nil
		}

		return models.Session{}, ape.ErrorInternal(appErr.Unwrap())
	}

	//Login flow
	session, appErr := a.sessions.Create(ctx, user, client)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, *ape.Error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	session, appErr := a.sessions.Refresh(ctx, sessionID, user, client, token)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) Logout(ctx context.Context, sessionID uuid.UUID) *ape.Error {
	_, appErr := a.GetSession(ctx, sessionID)
	if appErr != nil {
		return appErr
	}

	appErr = a.sessions.Delete(ctx, sessionID)
	if appErr != nil {
		return appErr
	}
	return nil
}
