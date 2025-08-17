package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/entity"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
)

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

func (a App) Login(ctx context.Context, email string, client string) (models.Session, models.TokensPair, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errx.ErrorUserNotFound) {
			err = a.users.Create(ctx, entity.UserCreateInput{
				Email:    email,
				Role:     roles.User,
				Verified: false,
			})
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

	session, tokensPair, err := a.sessions.Create(ctx, user, client)
	if err != nil {

		return models.Session{}, models.TokensPair{}, err
	}

	return session, tokensPair, nil
}

func (a App) GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (models.Session, error) {
	session, appErr := a.sessions.Get(ctx, sessionID, userID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	if session.UserID != userID {
		return models.Session{}, errx.RaiseSessionNotFound(
			ctx,
			fmt.Errorf("session %s does not belong to user %s", session.ID, userID),
			sessionID.String(),
			userID.String(),
		)
	}

	return session, nil
}

func (a App) GetUserSessions(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error) {
	_, appErr := a.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, pagination.Response{}, appErr
	}

	sessions, pagResp, appErr := a.sessions.SelectByUserID(ctx, userID, pag)
	if appErr != nil {
		return nil, pagination.Response{}, appErr
	}

	return sessions, pagResp, nil
}

func (a App) DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	session, appErr := a.GetUserSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return errx.RaiseSessionNotFound(
			ctx,
			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
			sessionID.String(),
			userID.String(),
		)
	}

	appErr = a.sessions.Delete(ctx, session.ID, userID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
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

// Note: if you want to check initiator rights, not in grpc-service package, with use middleware, you can uncomment the following lines
// and use the function AdminDeleteUserSessions in grpc-service package
//
//func (a App) AdminDeleteUserSessions(ctx context.Context, initiatorID, userID uuid.UUID) error {
//	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID)
//	if err != nil {
//		return err
//	}
//
//	appErr := a.sessions.Terminate(ctx, userID)
//	if appErr != nil {
//		return appErr
//	}
//	return nil
//}
//
// Note: if you want to check initiator rights, not in grpc-service package, with use middleware, you can uncomment the following lines
// and use the function AdminDeleteUserSession in grpc-service package
//
//func (a App) AdminDeleteUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error {
//	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID)
//	if err != nil {
//		return err
//	}
//
//	session, appErr := a.GetUserSession(ctx, userID, sessionID)
//	if appErr != nil {
//		return appErr
//	}
//
//	if session.UserID != userID {
//		return errx.RaiseSessionNotFound(
//			ctx,
//			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
//			sessionID.String(),
//			userID.String(),
//		)
//	}
//
//	appErr = a.sessions.Delete(ctx, session.ID, userID)
//	if appErr != nil {
//		return appErr
//	}
//
//	return nil
//}
