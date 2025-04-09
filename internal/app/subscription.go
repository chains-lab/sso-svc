package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app/ape"
	"github.com/hs-zavet/sso-oauth/internal/repo"
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
			return ape.ErrAccountNotFound
		default:
			return err
		}
	}

	return nil
}
