package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/chains-lab/sso-svc/internal/infra/password"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	oldPassword, newPassword string,
) error {
	user, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusBlocked {
		return errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user %s is blocked", userID),
		)
	}

	userRow, err := s.db.Users().FilterID(userID).Get(ctx)
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

	if err = bcrypt.CompareHashAndPassword([]byte(userRow.PasswordHash), []byte(oldPassword)); err != nil {
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

	now := time.Now().UTC()
	err = s.db.Users().Transaction(ctx, func(ctx context.Context) error {
		err = s.db.Users().FilterID(userID).UpdatePassword(hashStr, now).Update(ctx, now)
		if err != nil {
			return err
		}

		err = s.db.Sessions().FilterUserID(userID).Delete(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}
