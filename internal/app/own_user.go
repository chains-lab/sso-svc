package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
)

// DeleteOwnUser allows a user to delete their own account. deprecated: don't support in new design
func (a App) DeleteOwnUser(ctx context.Context, userID, initiatorSessionID uuid.UUID) error {
	_, err := a.getInitiator(ctx, userID, initiatorSessionID)
	if err != nil {
		return err
	}

	return a.users.Delete(ctx, userID)
}

func (a App) UpdatePassword(ctx context.Context, userID, SessionID uuid.UUID, currentPassword, newPassword string) error {
	user, err := a.getInitiator(ctx, userID, SessionID)
	if err != nil {
		return err
	}

	if err = a.users.CheckPassword(ctx, userID, currentPassword); err != nil {
		return errx.ErrorInvalidLogin.Raise(
			fmt.Errorf("current password is incorrect"),
		)
	}

	err = a.transaction(func(ctx context.Context) error {
		err = a.users.UpdatePassword(ctx, user.ID, newPassword)
		if err != nil {
			return err
		}

		err = a.session.DeleteAllForUser(ctx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
