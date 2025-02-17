package sqldb

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/service/data/dbx/sqldb/core"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/models"
)

type Accounts interface {
	Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error)
	GetByEmail(r *http.Request, email string) (*models.Account, error)
	GetByID(r *http.Request, id uuid.UUID) (*models.Account, error)
	UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error)
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

func (a *accounts) Create(r *http.Request, email string, role roles.UserRole) (*models.Account, error) {
	res, err := a.queries.CreateAccount(r.Context(), core.CreateAccountParams{
		Email: email,
		Role:  string(role),
	})
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) GetByID(r *http.Request, id uuid.UUID) (*models.Account, error) {
	res, err := a.queries.GetAccountByID(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) GetByEmail(r *http.Request, email string) (*models.Account, error) {
	res, err := a.queries.GetAccountByEmail(r.Context(), email)
	if err != nil {
		return nil, err
	}

	return parseAccount(res), nil
}

func (a *accounts) UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (*models.Account, error) {
	res, err := a.queries.UpdateAccountRole(r.Context(), core.UpdateAccountRoleParams{
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
