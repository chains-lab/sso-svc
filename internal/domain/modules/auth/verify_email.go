package auth

import (
	"context"
	"fmt"

	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/internal/domain/errx"
)

// VerifyEmail is a callback function to verify email addresses, which we should use in Kafka consumer
func (s Service) VerifyEmail(ctx context.Context, email string) (entity.AccountEmail, error) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.AccountEmail{}, err
	}

	emailData, err := s.db.UpdateAccountEmailVerification(ctx, account.ID, true)
	if err != nil {
		return entity.AccountEmail{}, errx.ErrorInternal.Raise(
			fmt.Errorf("verifying email for account %s, cause: %w", account.ID, err),
		)
	}

	return emailData, nil
}
