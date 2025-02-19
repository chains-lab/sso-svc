package sqldb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	core2 "github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb/core"
)

type Accounts interface {
	Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error)
}

type accounts struct {
	queries *core2.Queries
}

func NewAccounts(url string) (Accounts, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &accounts{queries: core2.New(db)}, nil
}

func (a *accounts) Create(ctx context.Context, email string, role roles.UserRole) (*models.Account, error) {
	acc, err := a.queries.CreateAccount(ctx, core2.CreateAccountParams{
		Email: email,
		Role:  string(role),
	})
	if err != nil {
		return nil, err
	}

	res, err := parseAccount(acc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *accounts) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	acc, err := a.queries.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res, err := parseAccount(acc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	acc, err := a.queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	res, err := parseAccount(acc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *accounts) UpdateRole(ctx context.Context, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	acc, err := a.queries.UpdateAccountRole(ctx, core2.UpdateAccountRoleParams{
		ID:   id,
		Role: string(role),
	})
	if err != nil {
		return nil, err
	}

	res, err := parseAccount(acc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func parseAccount(account core2.Account) (*models.Account, error) {
	role, err := roles.ParseUserRole(account.Role)
	if err != nil {
		return nil, err
	}

	return &models.Account{
		ID:        account.ID,
		Email:     account.Email,
		Role:      role,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
