package app

import (
	"context"
	"errors"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/ape"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	if user.Suspended {
		return models.Session{}, models.TokensPair{}, ape.RaiseUserSuspended(user.ID)
	}

	session, tokensPair, appErr := a.sessions.Refresh(ctx, sessionID, user, client, token)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	return session, tokensPair, nil
}

func (a App) Login(ctx context.Context, email string, role roles.Role, client string) (models.Session, models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ape.ErrUserNotFound) {
			err = a.users.Create(ctx, email, role)
			if err != nil {

				// If we fail to create a user, we should return an internal error
				return models.Session{}, models.TokensPair{}, err
			}

			user, err = a.users.GetByEmail(ctx, email)
			if err != nil {

				// It a good return internal error here anyway, because we already created the user in logic above
				return models.Session{}, models.TokensPair{}, err
			}

			session, tokensPair, err := a.sessions.Create(ctx, user, client)
			if err != nil {

				// If we fail to create a session, we should return an internal error
				return models.Session{}, models.TokensPair{}, err
			}

			return session, tokensPair, nil
		}

		return models.Session{}, models.TokensPair{}, err
	}

	if user.Suspended {

		return models.Session{}, models.TokensPair{}, ape.RaiseUserSuspended(user.ID)
	}

	session, tokensPair, err := a.sessions.Create(ctx, user, client)
	if err != nil {

		return models.Session{}, models.TokensPair{}, err
	}

	return session, tokensPair, nil
}

func (a App) GetSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, appErr := a.sessions.Get(ctx, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	if session.UserID != userID {
		return models.Session{}, ape.RaiseSessionDoesNotBelongToUser(sessionID, userID)
	}

	return session, nil
}

func (a App) GetUserSessions(ctx context.Context, userID uuid.UUID, page, limit uint64) ([]models.Session, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	if user.Suspended {
		return nil, ape.RaiseUserSuspended(user.ID)
	}

	sessions, appErr := a.sessions.SelectByUserID(ctx, userID, page, limit)
	if appErr != nil {
		return nil, appErr
	}

	return sessions, nil
}

func (a App) DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	session, appErr := a.GetSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return ape.RaiseSessionDoesNotBelongToUser(sessionID, userID)
	}

	appErr = a.sessions.Delete(ctx, session.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	user, appError := a.GetUserByID(ctx, userID)
	if appError != nil {
		return appError
	}

	if user.Suspended {
		return ape.RaiseUserSuspended(user.ID)
	}

	appErr := a.sessions.Terminate(ctx, userID)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a App) AdminDeleteSessions(ctx context.Context, initiatorID, userID uuid.UUID) error {
	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.RaiseNoPermissions(err)
		}
	}

	appErr := a.sessions.Terminate(ctx, userID)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a App) AdminDeleteSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error {
	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.RaiseNoPermissions(err)
		}
	}

	session, appErr := a.GetSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return ape.RaiseSessionDoesNotBelongToUser(session.ID, userID)
	}

	appErr = a.sessions.Delete(ctx, session.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) AdminDeleteUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error {
	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return err
	}

	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return ape.RaiseNoPermissions(err)
		}
	}

	session, appErr := a.GetSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return ape.RaiseSessionDoesNotBelongToUser(session.ID, userID)
	}

	appErr = a.sessions.Delete(ctx, session.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}
