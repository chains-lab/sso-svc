package sqldb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/repository/sqldb/core"
	"github.com/recovery-flow/tokens/identity"
)

type Accounts struct {
	queries *core.Queries
}

func NewAccounts(url string) (*Accounts, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Accounts{queries: core.New(db)}, nil
}

func (a *Accounts) Insert(ctx context.Context, email string, idn identity.IdnType) (*models.Account, error) {
	acc, err := a.queries.CreateAccount(ctx, core.CreateAccountParams{
		Email: email,
		Role:  string(idn),
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

func (a *Accounts) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
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

func (a *Accounts) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
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

func (a *Accounts) UpdateRole(ctx context.Context, id uuid.UUID, idn identity.IdnType) (*models.Account, error) {
	acc, err := a.queries.UpdateAccountRole(ctx, core.UpdateAccountRoleParams{
		ID:   id,
		Role: string(idn),
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

func parseAccount(account core.Account) (*models.Account, error) {
	role, err := identity.ParseIdentityType(account.Role)
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
