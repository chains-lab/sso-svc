package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/sso-svc/internal/domain/entity"
	"github.com/umisto/sso-svc/internal/domain/modules/auth"
	"github.com/umisto/sso-svc/internal/repo/pgdb"
)

func (r *Repository) CreateAccount(ctx context.Context, params auth.CreateAccountParams) (entity.Account, error) {
	var account entity.Account

	err := r.sql.accounts.Transaction(ctx, func(ctx context.Context) error {
		now := time.Now().UTC()
		accountID := uuid.New()

		acc := pgdb.Account{
			ID:                accountID,
			Username:          params.Username,
			Role:              params.Role,
			Status:            entity.AccountStatusActive,
			CreatedAt:         now,
			UpdatedAt:         now,
			UsernameUpdatedAt: now,
		}

		account = acc.ToEntity()

		err := r.sql.accounts.Insert(ctx, acc)
		if err != nil {
			return err
		}

		emailRow := pgdb.AccountEmail{
			AccountID: accountID,
			Email:     params.Email,
			Verified:  false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = r.sql.emails.Insert(ctx, emailRow)
		if err != nil {
			return err
		}

		passwordRow := pgdb.AccountPassword{
			AccountID: accountID,
			Hash:      params.PasswordHash,
			CreatedAt: now,
			UpdatedAt: now,
		}

		return r.sql.passwords.Insert(ctx, passwordRow)
	})
	if err != nil {
		return entity.Account{}, err
	}

	return account, err
}

func (r *Repository) GetAccountByID(ctx context.Context, accountID uuid.UUID) (entity.Account, error) {
	acc, err := r.sql.accounts.New().FilterID(accountID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToEntity(), nil
}

func (r *Repository) GetAccountByUsername(ctx context.Context, username string) (entity.Account, error) {
	acc, err := r.sql.accounts.New().FilterUsername(username).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToEntity(), nil
}

func (r *Repository) GetAccountByEmail(ctx context.Context, email string) (entity.Account, error) {
	acc, err := r.sql.accounts.New().FilterEmail(email).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToEntity(), nil
}

func (r *Repository) UpdateAccountUsername(ctx context.Context, accountID uuid.UUID, newUsername string) (entity.Account, error) {
	var account entity.Account

	err := r.sql.accounts.Transaction(ctx, func(ctx context.Context) error {
		accs, err := r.sql.accounts.New().
			FilterID(accountID).
			UpdateUsername(newUsername, time.Now().UTC()).
			Update(ctx)
		if err != nil {
			return err
		}

		if len(accs) == 1 {
			account = accs[0].ToEntity()
		} else {
			return fmt.Errorf("expected to update 1 account, updated %d", len(accs))
		}

		err = r.DeleteSessionsForAccount(ctx, accountID)
		if err != nil {
			return err
		}

		return nil
	})

	return account, err
}

func (r *Repository) UpdateAccountStatus(ctx context.Context, accountID uuid.UUID, status string) (entity.Account, error) {
	accs, err := r.sql.accounts.New().
		FilterID(accountID).
		UpdateStatus(status).
		Update(ctx)
	if err != nil {
		return entity.Account{}, err
	}

	if len(accs) != 1 {
		return entity.Account{}, fmt.Errorf("expected to update 1 account, updated %d", len(accs))
	}
	return accs[0].ToEntity(), nil
}

func (r *Repository) GetAccountEmail(ctx context.Context, accountID uuid.UUID) (entity.AccountEmail, error) {
	acc, err := r.sql.emails.New().FilterAccountID(accountID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.AccountEmail{}, nil
	case err != nil:
		return entity.AccountEmail{}, err
	}

	return acc.ToEntity(), nil
}

func (r *Repository) UpdateAccountEmailVerification(
	ctx context.Context,
	accountID uuid.UUID,
	verified bool,
) (entity.AccountEmail, error) {
	accs, err := r.sql.emails.New().
		FilterAccountID(accountID).
		UpdateVerified(verified).
		Update(ctx)
	if err != nil {
		return entity.AccountEmail{}, err
	}

	if len(accs) != 1 {
		return entity.AccountEmail{}, fmt.Errorf("expected to update 1 account, updated %d", len(accs))
	}
	return accs[0].ToEntity(), nil
}

func (r *Repository) GetAccountPassword(ctx context.Context, accountID uuid.UUID) (entity.AccountPassword, error) {
	acc, err := r.sql.passwords.New().FilterAccountID(accountID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.AccountPassword{}, nil
	case err != nil:
		return entity.AccountPassword{}, err
	}

	return acc.ToEntity(), nil
}

func (r *Repository) UpdateAccountPassword(
	ctx context.Context,
	accountID uuid.UUID,
	passwordHash string,
) (entity.AccountPassword, error) {
	var password entity.AccountPassword

	err := r.sql.accounts.Transaction(ctx, func(ctx context.Context) error {
		accs, err := r.sql.passwords.New().
			FilterAccountID(accountID).
			UpdateHash(passwordHash).
			Update(ctx)
		if err != nil {
			return err
		}

		if len(accs) != 1 {
			return fmt.Errorf("expected to update 1 account, updated %d", len(accs))
		}

		password = accs[0].ToEntity()

		return r.DeleteSessionsForAccount(ctx, accountID)
	})
	if err != nil {
		return entity.AccountPassword{}, err
	}

	return password, nil
}

func (r *Repository) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	return r.sql.accounts.New().FilterID(accountID).Delete(ctx)
}
