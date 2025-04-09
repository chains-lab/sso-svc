package app

import (
	"context"
	"time"

	"github.com/google/uuid"
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
		return err
	}

	return nil
}
