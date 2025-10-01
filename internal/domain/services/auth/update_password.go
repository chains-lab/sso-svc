package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/sso-svc/internal/domain/errx"
	"github.com/chains-lab/sso-svc/internal/domain/models"
	"github.com/chains-lab/sso-svc/internal/infra/password"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	oldPassword, newPassword string,
) error {
	initiator, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get initiator with id '%s', cause: %w", userID, err),
		)
	}

	if initiator == (models.User{}) {
		return errx.ErrorUnauthenticated.Raise(
			fmt.Errorf("initiator with id '%s' not found", userID),
		)
	}

	if initiator.Status == enum.UserStatusBlocked {
		return errx.ErrorInitiatorIsBlocked.Raise(
			fmt.Errorf("user %s is blocked", userID),
		)
	}

	passData, err := s.db.GetUserPassword(ctx, userID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("getting password for user %s, cause: %w", userID, err),
		)
	}

	if (passData == models.UserPassword{}) {
		return errx.ErrorUserNotFound.Raise(
			fmt.Errorf("password for user %s not found, cause: %w", userID, err),
		)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passData.Hash), []byte(oldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errx.ErrorInvalidLogin.Raise(
				fmt.Errorf("invalid credentials for user %s, cause: %w", userID, err),
			)
		}

		return errx.ErrorInternal.Raise(
			fmt.Errorf("comparing newPassword hash for user %s, cause: %w", userID, err),
		)
	}

	err = password.ReliabilityCheck(newPassword)
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

	err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.UpdateUserPassword(ctx, userID, hashStr, time.Now().UTC())
		if err != nil {
			return err
		}

		err = s.db.DeleteAllSessionsForUser(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}
