package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u Service) CheckPassword(ctx context.Context, userID uuid.UUID, password string) error {
	secret, err := u.db.UsersPassword().FilterID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorUserNotFound.Raise(
				fmt.Errorf("password for user %s not found, cause: %w", userID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("getting password for user %s, cause: %w", userID, err),
			)
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials for user %s, cause: %w", userID, err),
			)
		}

		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing password hash for user %s, cause: %w", userID, err),
		)
	}

	return nil
}
