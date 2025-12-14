package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chains-lab/sso-svc/internal/domain"
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

func (r *Repository) CreateAccount(ctx context.Context, params domain.CreateAccountParams) (entity.Account, error) {
	acc, err := r.sql.CreateAccount(ctx, pgdb.CreateAccountParams{
		Username: params.Username,
		Role:     pgdb.AccountRole(params.Role),
		Status:   pgdb.AccountStatusActive,
	})
	if err != nil {
		return entity.Account{}, err
	}

	_, err = r.sql.CreateAccountEmail(ctx, pgdb.CreateAccountEmailParams{
		AccountID: acc.ID,
		Email:     params.Email,
		Verified:  false,
	})
	if err != nil {
		return entity.Account{}, err
	}

	_, err = r.sql.CreateAccountPassword(ctx, pgdb.CreateAccountPasswordParams{
		AccountID: acc.ID,
		Hash:      params.PasswordHash,
	})
	if err != nil {
		return entity.Account{}, err
	}

	return acc.ToModel(), err
}

func (r *Repository) GetAccountByID(ctx context.Context, accountID uuid.UUID) (entity.Account, error) {
	acc, err := r.sql.GetAccountByID(ctx, accountID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToModel(), nil
}

func (r *Repository) GetAccountByUsername(ctx context.Context, username string) (entity.Account, error) {
	acc, err := r.sql.GetAccountByUsername(ctx, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToModel(), nil
}

func (r *Repository) GetAccountByEmail(ctx context.Context, email string) (entity.Account, error) {
	acc, err := r.sql.GetAccountByEmail(ctx, email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Account{}, nil
	case err != nil:
		return entity.Account{}, err
	}

	return acc.ToModel(), nil
}

func (r *Repository) UpdateAccountUsername(
	ctx context.Context,
	accountID uuid.UUID,
	newUsername string,
) (entity.Account, error) {
	var account entity.Account

	err := r.Transaction(ctx, func(ctx context.Context) error {
		row, err := r.sql.UpdateAccountUsername(ctx, pgdb.UpdateAccountUsernameParams{
			ID:       accountID,
			Username: newUsername,
		})
		if err != nil {
			return err
		}

		err = r.DeleteSessionsForAccount(ctx, accountID)
		if err != nil {
			return err
		}

		account = row.ToModel()

		return nil
	})

	return account, err
}

func (r *Repository) UpdateAccountStatus(
	ctx context.Context,
	accountID uuid.UUID,
	status string,
) (entity.Account, error) {
	acc, err := r.sql.UpdateAccountStatus(ctx, pgdb.UpdateAccountStatusParams{
		ID:     accountID,
		Status: pgdb.AccountStatus(status),
	})
	if err != nil {
		return entity.Account{}, err
	}

	return acc.ToModel(), nil
}

func (r *Repository) GetAccountEmail(ctx context.Context, accountID uuid.UUID) (entity.AccountEmail, error) {
	email, err := r.sql.GetAccountEmail(ctx, accountID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.AccountEmail{}, nil
	case err != nil:
		return entity.AccountEmail{}, err
	}

	return email.ToModel(), nil
}

func (r *Repository) UpdateAccountEmailVerification(
	ctx context.Context,
	accountID uuid.UUID,
	verified bool,
) (entity.AccountEmail, error) {
	email, err := r.sql.UpdateVerifiedEmail(ctx, pgdb.UpdateVerifiedEmailParams{
		AccountID: accountID,
		Verified:  verified,
	})
	if err != nil {
		return entity.AccountEmail{}, err
	}

	return email.ToModel(), nil
}

func (r *Repository) GetAccountPassword(ctx context.Context, accountID uuid.UUID) (entity.AccountPassword, error) {
	u, err := r.sql.GetAccountPassword(ctx, accountID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.AccountPassword{}, nil
	case err != nil:
		return entity.AccountPassword{}, err
	}

	return u.ToModel(), nil
}

func (r *Repository) UpdateAccountPassword(
	ctx context.Context,
	accountID uuid.UUID,
	passwordHash string,
) error {
	return r.Transaction(ctx, func(ctx context.Context) error {
		_, err := r.sql.UpdateAccountPassword(ctx, pgdb.UpdateAccountPasswordParams{
			AccountID: accountID,
			Hash:      passwordHash,
		})
		if err != nil {
			return err
		}

		err = r.DeleteSessionsForAccount(ctx, accountID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	err := r.sql.DeleteAccount(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}
