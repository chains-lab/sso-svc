package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/chains-auth/internal/app/ape"
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/chains-auth/internal/jwtkit"
	"github.com/chains-lab/chains-auth/internal/repo"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type accountsRepo interface {
	Create(ctx context.Context, input repo.AccountCreateRequest) error
	Update(ctx context.Context, ID uuid.UUID, input repo.AccountUpdateRequest) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (repo.Account, error)
	GetByEmail(ctx context.Context, email string) (repo.Account, error)
	Transaction(fn func(ctx context.Context) error) error
	Drop(ctx context.Context) error
}

type Accounts struct {
	repo accountsRepo
	jwt  JWTManager
}

func NewAccount(cfg config.Config, log *logrus.Logger) (Accounts, error) {
	data, err := repo.NewAccounts(cfg, log)
	if err != nil {
		return Accounts{}, nil
	}

	jwt := jwtkit.NewManager(cfg)

	return Accounts{
		repo: data,
		jwt:  jwt,
	}, nil
}

func (a Accounts) Create(ctx context.Context, email string, role roles.Role) *ape.Error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.repo.Create(ctx, repo.AccountCreateRequest{
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
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (a Accounts) UpdateRole(ctx context.Context, ID uuid.UUID, role roles.Role) *ape.Error {
	UpdatedAt := time.Now().UTC()

	err := a.repo.Update(ctx, ID, repo.AccountUpdateRequest{
		Role:      &role,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorAccountDoesNotExist(ID, err)
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (a Accounts) GetByID(ctx context.Context, ID uuid.UUID) (models.Account, *ape.Error) {
	account, err := a.repo.GetByID(ctx, ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Account{}, ape.ErrorAccountDoesNotExist(ID, err)
		default:
			return models.Account{}, ape.ErrorInternal(err)
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

func (a Accounts) GetByEmail(ctx context.Context, email string) (models.Account, *ape.Error) {
	account, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Account{}, ape.ErrorAccountDoesNotExistByEmail(email, err)
		default:
			return models.Account{}, ape.ErrorInternal(err)
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
