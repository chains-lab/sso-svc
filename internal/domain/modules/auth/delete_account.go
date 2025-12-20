package auth

import (
	"context"
	"fmt"

	"github.com/umisto/sso-svc/internal/domain/errx"
)

func (s Service) DeleteOwnAccount(ctx context.Context, initiator InitiatorData) error {
	_, _, err := s.ValidateSession(ctx, initiator)
	if err != nil {
		return err
	}

	err = s.db.DeleteAccount(ctx, initiator.AccountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete account with id: %s, cause: %w", initiator.AccountID, err),
		)
	}

	return nil
}
