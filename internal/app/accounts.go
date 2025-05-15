package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func (a App) CreateAccount(ctx context.Context, email string, role roles.Role) *ape.Error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.accounts.Create(ctx, repo.AccountCreateRequest{
		ID:           ID,
		Email:        email,
		Role:         role,
		Subscription: uuid.Nil,
		CreatedAt:    CreatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorAccountAlreadyExists(err)
		default:
			return ape.ErrorInternalServer(err)
		}
	}

	return nil
}

func (a App) UpdateAccountRole(ctx context.Context, ID uuid.UUID, role, initiatorRole roles.Role) *ape.Error {
	UpdatedAt := time.Now().UTC()

	if roles.CompareRolesUser(role, initiatorRole) != 1 {
		return ape.ErrorUserNoPermissionToUpdateRole(fmt.Errorf("user has no permission to update role"))
	}

	err := a.accounts.Update(ctx, ID, repo.AccountUpdateRequest{
		Role:      &role,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorAccountDoesNotExistByID(ID, err)
		default:
			return ape.ErrorInternalServer(err)
		}
	}

	return nil
}

func (a App) GetAccountByID(ctx context.Context, ID uuid.UUID) (models.Account, *ape.Error) {
	account, err := a.accounts.GetByID(ctx, ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Account{}, ape.ErrorAccountDoesNotExistByID(ID, err)
		default:
			return models.Account{}, ape.ErrorInternalServer(err)
		}
	}

	return models.Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         account.Role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}, nil
}

func (a App) GetAccountByEmail(ctx context.Context, email string) (models.Account, *ape.Error) {
	account, err := a.accounts.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Account{}, ape.ErrorAccountDoesNotExistByEmail(email, err)
		default:
			return models.Account{}, ape.ErrorInternalServer(err)
		}
	}

	return models.Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         account.Role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}, nil
}
