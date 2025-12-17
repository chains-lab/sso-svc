package auth

import (
	"context"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
)

func (s Service) UpdateUsername(
	ctx context.Context,
	initiator InitiatorData,
	password string,
	newUsername string,
) (entity.Account, error) {
	account, _, err := s.ValidateSession(ctx, initiator)
	if err != nil {
		return entity.Account{}, err
	}

	//if err = account.CanChangeUsername(); err != nil {
	//	return entity.Account{}, err
	//}

	if err = s.CheckUsernameRequirements(newUsername); err != nil {
		return entity.Account{}, err
	}

	if err = s.checkAccountPassword(ctx, initiator.AccountID, password); err != nil {
		return entity.Account{}, err
	}

	account, err = s.db.UpdateAccountUsername(ctx, initiator.AccountID, newUsername)
	if err != nil {
		return entity.Account{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating username for account %s, cause: %w", initiator.AccountID, err),
		)
	}

	email, err := s.GetAccountEmail(ctx, account.ID)
	if err != nil {
		return entity.Account{}, err
	}

	err = s.event.WriteAccountUsernameChanged(ctx, account, email.Email)
	if err != nil {
		return entity.Account{}, err
	}

	return account, nil
}
