package domain

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	password string,
	newUsername string,
) (entity.Account, error) {
	account, err := s.GetInitiator(ctx, accountID)
	if err != nil {
		return entity.Account{}, err
	}

	if account.CanChangeUsername() != nil {
		return entity.Account{}, errx.ErrorCannotChangeUsernameYet.Raise(
			fmt.Errorf("account %s cannot change username due to cooldown, cause: %w", accountID, err),
		)
	}

	if err = s.CheckUsernameRequirements(newUsername); err != nil {
		return entity.Account{}, err
	}

	if err = s.CheckAccountPassword(ctx, accountID, password); err != nil {
		return entity.Account{}, err
	}

	account, err = s.db.UpdateAccountUsername(ctx, accountID, newUsername)
	if err != nil {
		return entity.Account{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating username for account %s, cause: %w", accountID, err),
		)
	}
	if account.IsNil() {
		return entity.Account{}, errx.ErrorInitiatorNotFound.Raise(
			fmt.Errorf("account %s not found when updating username, cause: %w", accountID, err),
		)
	}

	return account, nil
}
