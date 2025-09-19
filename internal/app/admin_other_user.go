package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) AdminGetUser(ctx context.Context, initiatorID, initiatorSessionID, userID uuid.UUID) (models.User, error) {
	_, user, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
	if err != nil {
		return models.User{}, err
	}

	_, err = a.session.GetSessionForInitiator(ctx, initiatorID, initiatorSessionID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a App) AdminDeleteUser(
	ctx context.Context,
	initiatorID, initiatorSessionID, userID uuid.UUID,
) error {
	_, _, err := a.users.CompareRightsForAdmins(ctx, initiatorID, userID, 1)
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
