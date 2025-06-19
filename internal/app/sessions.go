package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type sessionsDomain interface {
	Terminate(ctx context.Context, userUD uuid.UUID) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
	Get(ctx context.Context, sessionID uuid.UUID) (models.Session, error)
	SelectByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
	Create(ctx context.Context, user models.User, client string) (models.Session, models.TokensPair, error)
	Refresh(ctx context.Context, sessionID uuid.UUID, user models.User, client, token string) (models.Session, models.TokensPair, error)
}

func (a App) Refresh(ctx context.Context, userID, sessionID uuid.UUID, client, token string) (models.Session, models.TokensPair, error) {
	user, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	session, tokensPair, appErr := a.sessions.Refresh(ctx, sessionID, user, client, token)
	if appErr != nil {
		return models.Session{}, models.TokensPair{}, appErr
	}

	return session, tokensPair, nil
}

func (a App) Login(ctx context.Context, email string, role roles.Role, client string) (models.Session, models.TokensPair, error) {
	user, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ape.ErrUserDoesNotExist) {
			err = a.users.Create(ctx, email, role)
			if err != nil {

				// If we fail to create a user, we should return an internal error
				return models.Session{}, models.TokensPair{}, ape.ErrorInternal(fmt.Errorf("error in login %s", err))
			}

			user, err = a.users.GetByEmail(ctx, email)
			if err != nil {

				// It a good return internal error here anyway, because we already created the user in logic above
				return models.Session{}, models.TokensPair{}, ape.ErrorInternal(fmt.Errorf("error in login %s", err))
			}

			session, tokensPair, err := a.sessions.Create(ctx, user, client)
			if err != nil {

				// If we fail to create a session, we should return an internal error
				return models.Session{}, models.TokensPair{}, ape.ErrorInternal(fmt.Errorf("error in login %s", err))
			}

			return session, tokensPair, nil
		}

		return models.Session{}, models.TokensPair{}, ape.ErrorInternal(fmt.Errorf("error in login %s", err))
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
		return models.Session{}, ape.ErrorSessionDoesNotBelongToUser(sessionID, userID)
	}

	return session, nil
}

func (a App) SelectSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	_, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	sessions, appErr := a.sessions.SelectByUserID(ctx, userID)
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
		return ape.ErrorSessionDoesNotBelongToUser(sessionID, userID)
	}

	appErr = a.sessions.Delete(ctx, session.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) TerminateSessions(ctx context.Context, userID uuid.UUID) error {
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

func (a App) AdminTerminateSessions(ctx context.Context, initiatorID, userID uuid.UUID) error {
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
			return ape.ErrNoPermission
		}
	}

	appErr := a.sessions.Terminate(ctx, userID)
	if appErr != nil {
		return appErr
	}
	return nil
}

//func (a App) GetSession(ctx context.Context, sessionID uuid.UUID) (models.Session, error) {
//	session, appErr := a.sessions.Get(ctx, sessionID)
//	if appErr != nil {
//		return models.Session{}, appErr
//	}
//
//	return session, nil
//}

//func (a App) AdminTerminateSessions(ctx context.Context, userID uuid.UUID) error {
//	user, err := a.GetUserByID(ctx, userID)
//	if err != nil {
//		return err
//	}
//
//	if user.Role == roles.SuperUser {
//		return ape.ErrNoPermission
//	}
//
//	appErr := a.sessions.Terminate(ctx, userID)
//	if appErr != nil {
//		return appErr
//	}
//	return nil
//}
