package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) UpdatePassword(
	ctx context.Context,
	accountID uuid.UUID,
	oldPassword, newPassword string,
) error {
	_, err := s.GetInitiator(ctx, accountID)
	if err != nil {
		return err
	}

	passData, err := s.db.GetAccountPassword(ctx, accountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("getting password for account %s, cause: %w", accountID, err),
		)
	}
	if passData.IsNil() {
		return errx.ErrorInitiatorNotFound.Raise(
			fmt.Errorf("password for account %s not found, cause: %w", accountID, err),
		)
	}
	if err = passData.CanChangePassword(); err != nil {
		return err
	}

	if err = s.CheckAccountPassword(ctx, accountID, oldPassword); err != nil {
		return err
	}

	if err = s.CheckPasswordRequirements(newPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing new newPassword for account '%s', cause: %w", accountID, err),
		)
	}

	err = s.db.UpdateAccountPassword(ctx, accountID, string(hash))
	if err != nil {
		return err
	}

	return nil
}
