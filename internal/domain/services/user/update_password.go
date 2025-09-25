package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/services/user/password"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u Service) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	oldPassword, newPassword string,
) error {
	user, err := u.GetInitiator(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusBlocked {
		return errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user %s is blocked", userID),
		)
	}

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

	if err = bcrypt.CompareHashAndPassword([]byte(secret.PassHash), []byte(oldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials for user %s, cause: %w", userID, err),
			)
		}

		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing newPassword hash for user %s, cause: %w", userID, err),
		)
	}

	err = password.CheckPassword(newPassword)
	if err != nil {
		return errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing new newPassword for user '%s', cause: %w", userID, err),
		)
	}

	hashStr := string(hash)

	err = u.db.UsersPassword().FilterID(userID).Update(ctx, schemas.UserPassUpdateInput{
		PassHash:  &hashStr,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("updating newPassword for user '%s', cause: %w", userID, err),
		)
	}

	err = u.db.Transaction(ctx, func(ctx context.Context) error {
		err = u.db.UsersPassword().FilterID(userID).Update(ctx, schemas.UserPassUpdateInput{
			PassHash:  &hashStr,
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		err = u.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}
