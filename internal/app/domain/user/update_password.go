package user

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u User) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	err := checkPassword(password)
	if err != nil {
		return errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing new password for user '%s', cause: %w", userID, err),
		)
	}

	err = u.passQ.New().FilterID(userID).Update(ctx,
		map[string]interface{}{
			"password_hash": string(hash),
			"updated_at":    time.Now().UTC(),
		},
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating password for user '%s', cause: %w", userID, err),
		)
	}

	return nil
}
