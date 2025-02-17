package sqldb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service/data/sqldb/core"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
)

type Accounts interface {
	Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error)
}

type accounts struct {
	queries *core.Queries
}

func NewAccounts(url string) (Accounts, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &accounts{queries: core.New(db)}, nil
}

func (a *accounts) Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error) {
	res, err := a.queries.CreateAccount(ctx, core.CreateAccountParams{
		Email: email,
		Role:  string(role),
	})
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	res, err := a.queries.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	res, err := a.queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	res, err := a.queries.UpdateAccountRole(ctx, core.UpdateAccountRoleParams{
		ID:   id,
		Role: string(role),
	})
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func parseAccount(account core.Account) *models.Account {
	return &models.Account{
		ID:        account.ID,
		Email:     account.Email,
		Role:      account.Role,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
