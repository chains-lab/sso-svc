package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a App) RegisterAdmin(ctx context.Context, initiatorID uuid.UUID, email, password, role string) (models.User, error) {
	user, err := a.GetUserByEmail(ctx, email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := a.GetUserByID(ctx, initiatorID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Role == roles.User || initiator.Role == roles.Moder {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, role) < 1 {
			return models.User{}, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to create user Role %s", initiator.Role, role),
			)
		}
	}

	txErr := a.usersQ.Transaction(func(ctx context.Context) error {
		err = a.usersQ.New().Insert(ctx, dbx.UserModel{
			ID:             uuid.New(),
			Email:          email,
			Role:           role,
			EmailVer:       true,
			EmailUpdatedAt: time.Now().UTC(),
			CreatedAt:      time.Now().UTC(),
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("creating user with email '%s': %w", email, err),
			)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("hashing password for user: %w", err),
			)
		}

		err = a.passQ.New().Insert(ctx, dbx.UserPasswordModel{
			ID:        uuid.New(),
			PassHash:  string(hash),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("creating user with password '%s': %w", password, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err = a.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// comparisonRightsForAdmins compares the roles of the initiator and the user to determine
// if the initiator has the necessary permissions to perform actions on the user.
//
// dif its - 1 if initiator role must be strictly greater than user role, 0 - equal, -1 - less
func (a App) comparisonRightsForAdmins(ctx context.Context, initiatorID, userID uuid.UUID, dif int) (initiator, user models.User, err error) {
	initiator, err = a.GetInitiatorByID(ctx, initiatorID)
	if err != nil {
		return initiator, user, err
	}

	user, err = a.GetUserByID(ctx, userID)
	if err != nil {
		return initiator, user, err
	}

	if user.Role == roles.User {
		return initiator, user, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("user %s is already a user", userID),
		)
	}

	if initiatorID == userID {
		return initiator, user, nil
	}

	if user.Role != roles.SuperUser {
		if roles.CompareRolesUser(initiator.Role, user.Role) < dif {
			return initiator, user, errx.ErrorNoPermissions.Raise(
				fmt.Errorf("initiator Role %s is not allowed to interact with this user", initiator.Role),
			)
		}
	}

	return initiator, user, nil
}

func (a App) AdminDeleteUserSessions(ctx context.Context, initiatorID, userID uuid.UUID) error {
	_, _, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 1)
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
	_, _, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	session, appErr := a.GetUserSession(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
		)
	}

	appErr = a.DeleteUserSession(ctx, session.ID, userID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) AdminGetUser(ctx context.Context, initiatorID, userID uuid.UUID) (models.User, error) {
	initiator, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 0)
	if err != nil {
		return models.User{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator %s does not have permission to access user %s", initiator.ID, userID),
		)
	}

	return user, nil
}

func (a App) AdminGetUserSession(ctx context.Context, initiatorID, userID, sessionID uuid.UUID) (models.Session, error) {
	initiator, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 0)
	if err != nil {
		return models.Session{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return models.Session{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator %s does not have permission to access session of user %s", initiator.ID, userID),
		)
	}

	session, appErr := a.GetUserSession(ctx, userID, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) AdminGetUserSessions(
	ctx context.Context,
	initiatorID, userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	initiator, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 0)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	if !(roles.CompareRolesUser(initiator.Role, user.Role) > 1 || initiator.Role == roles.SuperUser) {
		return nil, pagi.Response{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator %s does not have permission to access sessions of user %s", initiator.ID, userID),
		)
	}

	return a.GetUserSessions(ctx, userID, pag, sort)
}

func (a App) AdminDeleteAdmin(
	ctx context.Context,
	initiatorID, userID uuid.UUID,
) error {
	_, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	txErr := a.usersQ.New().Transaction(func(ctx context.Context) error {
		err := a.DeleteUserSessions(ctx, user.ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", user.ID, err),
			)
		}

		err = a.passQ.New().FilterID(user.ID).Delete(ctx)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting password for user %s: %w", user.ID, err),
			)
		}

		err = a.usersQ.New().FilterID(user.ID).Delete(ctx)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", user.ID, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (a App) AdminBlockUser(
	ctx context.Context,
	initiatorID, userID uuid.UUID,
) (models.User, error) {
	_, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	txErr := a.usersQ.New().Transaction(func(ctx context.Context) error {
		err := a.usersQ.New().FilterID(user.ID).Update(ctx,
			map[string]interface{}{"status": constant.UserStatusBlocked},
		)
		if err != nil {
			return errx.ErrorInternal.Raise(err)
		}

		err = a.DeleteUserSessions(ctx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (a App) DeleteAdmin(
	ctx context.Context,
	initiatorID, userID uuid.UUID,
) error {
	_, user, err := a.comparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	txErr := a.usersQ.New().Transaction(func(ctx context.Context) error {
	})

}
