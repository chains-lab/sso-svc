package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/google/uuid"
)

func (a App) SubscriptionUpdate(ctx context.Context, AccountID uuid.UUID, subscriptionID uuid.UUID) error {
	if err := a.sessions.Terminate(ctx, AccountID); err != nil {
		return err
	}

	if err := a.accounts.Update(ctx, AccountID, repo.AccountUpdateRequest{
		Subscription: &subscriptionID,
		UpdatedAt:    time.Now().UTC(),
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrAccountDoseNotExits
		default:
			return err
		}
	}

	return nil
}
