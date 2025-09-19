package user

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/dbx"
	"github.com/chains-lab/sso-svc/internal/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u User) Create(ctx context.Context, email, password, role string) error {
	err := roles.ParseRole(role)
	if err != nil {
		return errx.ErrorRoleNotSupported.Raise(
			fmt.Errorf("parsing role for new user with email '%s', cause: %w", email, err),
		)
	}

	err = checkPassword(password)
	if err != nil {
		return errx.ErrorPasswordIsInappropriate.Raise(err)
	}

	id := uuid.New()
	now := time.Now().UTC()

	err = u.query.New().Insert(ctx, dbx.UserModel{
		ID:        id,
		Role:      role,
		Status:    enum.UserStatusActive,
		CreatedAt: now,
	})

	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting new user with email '%s', cause: %w", email, err),
		)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("hashing password for user '%s', cause: %w", id, err),
		)
	}

	err = u.passQ.New().Insert(ctx, dbx.UserPasswordModel{
		ID:        id,
		PassHash:  string(hash),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting password for new user with email '%s', cause: %w", email, err),
		)
	}

	err = u.emailQ.New().Insert(ctx, dbx.UserEmailModel{
		ID:       id,
		Email:    email,
		Verified: true,
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("inserting email for new user with email '%s', cause: %w", email, err),
		)
	}

	return nil
}
