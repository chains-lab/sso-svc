package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/pagination"
	"github.com/google/uuid"
)

func (a App) ComparisonRightsForAdmins(ctx context.Context, initiatorID, userID uuid.UUID) (initiator, user models.User, err error) {
	initiator, err = a.GetInitiatorByID(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = a.GetUserByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if user.Role == roles.User {
		return initiator, user, errx.RaiseNoPermissions(
			ctx,
			fmt.Errorf("initiator Role %s is not allowed to interact wit with this user", initiator.Role),
		)
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < 1 {
			return initiator, user, errx.RaiseNoPermissions(
				ctx,
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}

func (a App) AdminDeleteUserSessions(ctx context.Context, initiatorID, userID uuid.UUID) error {
	_, _, err := a.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	if err != nil {
		return err
	}

	appErr := a.DeleteUserSessions(ctx, userID)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a App) AdminDeleteUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) error {
	_, _, err := a.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	if err != nil {
		return err
	}

	session, appErr := a.GetUserSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return errx.RaiseSessionNotFound(
			ctx,
			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
			sessionID,
			userID,
		)
	}

	appErr = a.DeleteUserSession(ctx, session.ID, userID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) AdminGetUser(ctx context.Context, initiatorID, userID uuid.UUID) (models.User, error) {
	initiator, user, err := a.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	if err != nil {
		return models.User{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return models.User{}, errx.RaiseNoPermissions(
			ctx,
			fmt.Errorf("initiator %s does not have permission to access user %s", initiator.ID, userID),
		)
	}

	return user, nil
}

func (a App) AdminGetUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) (models.Session, error) {
	initiator, user, err := a.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	if err != nil {
		return models.Session{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return models.Session{}, errx.RaiseNoPermissions(
			ctx,
			fmt.Errorf("initiator %s does not have permission to access session of user %s", initiator.ID, userID),
		)
	}

	session, appErr := a.GetUserSession(ctx, userID, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) AdminGetUserSessions(ctx context.Context, initiatorID, userID uuid.UUID, pag pagination.Request) ([]models.Session, pagination.Response, error) {
	initiator, user, err := a.ComparisonRightsForAdmins(ctx, initiatorID, userID)
	if err != nil {
		return nil, pagination.Response{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return nil, pagination.Response{}, errx.RaiseNoPermissions(
			ctx,
			fmt.Errorf("initiator %s does not have permission to access sessions of user %s", initiator.ID, userID),
		)
	}

	return a.GetUserSessions(ctx, userID, pag)
}
