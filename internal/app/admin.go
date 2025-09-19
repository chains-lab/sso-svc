package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) RegisterAdmin(ctx context.Context, initiatorID, initiatorSessionID uuid.UUID, email, password, role string) (models.User, error) {
	_, err := a.users.GetByEmail(ctx, email)
	if !errors.Is(err, errx.ErrorUserNotFound) {
		return models.User{}, err
	}

	initiator, err := a.getInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	if initiator.Role == roles.User || initiator.Role == roles.Moder {
		return models.User{}, errx.ErrorNoPermissions.Raise(
			fmt.Errorf("initiator with role %s is not allowed to create user", initiator.Role),
		)
	}

	err = a.users.Create(ctx, email, password, role)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user by email %s after creation: %w", email, err),
		)
	}

	return user, nil
}

func (a App) AdminDeleteUserSessions(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) error {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return err
	}

	err = a.transaction(func(ctx context.Context) error {
		err = a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (a App) AdminDeleteUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) error {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return err
	}

	session, appErr := a.session.GetForUser(ctx, userID, sessionID)
	if appErr != nil {
		return appErr
	}

	if session.UserID != userID {
		return errx.ErrorSessionNotFound.Raise(
			fmt.Errorf("session %s does not belong to user %s", sessionID, userID),
		)
	}

	appErr = a.session.DeleteOneForUser(ctx, session.ID, userID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a App) AdminGetUser(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) (models.User, error) {
	_, user, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a App) AdminGetUserSession(ctx context.Context, initiatorID, initiatorSessionID, userID, sessionID uuid.UUID) (models.Session, error) {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.Session{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.Session{}, err
	}

	session, appErr := a.session.GetForUser(ctx, userID, sessionID)
	if appErr != nil {
		return models.Session{}, appErr
	}

	return session, nil
}

func (a App) AdminSelectUserSessions(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Session, pagi.Response, error) {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	return a.session.SelectForUSer(ctx, userID, pag, sort)
}

func (a App) AdminDeleteUser(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
) error {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return err
	}

	txErr := a.transaction(func(ctx context.Context) error {
		err := a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", userID, err),
			)
		}

		err = a.users.Delete(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", userID, err),
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
	initiatorID, initiatorSessionID, userID uuid.UUID,
) (models.User, error) {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	txErr := a.transaction(func(ctx context.Context) error {
		err := a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", userID, err),
			)
		}

		err = a.users.SetStatus(ctx, userID, enum.UserStatusBlocked)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", userID, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}

func (a App) AdminUnblockUser(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
) (models.User, error) {
	_, _, err := a.users.ComparisonRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	txErr := a.transaction(func(ctx context.Context) error {
		err := a.session.DeleteAllForUser(ctx, userID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting sessions for user %s: %w", userID, err),
			)
		}

		err = a.users.SetStatus(ctx, userID, enum.UserStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("deleting user %s: %w", userID, err),
			)
		}

		return nil
	})

	if txErr != nil {
		return models.User{}, txErr
	}

	user, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, txErr
	}

	return user, nil
}
