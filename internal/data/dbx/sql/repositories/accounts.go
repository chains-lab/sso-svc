package repositories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/roles"
	sqlcore2 "github.com/recovery-flow/sso-oauth/internal/data/dbx/sql/repositories/sqlcore"
)

type Accounts interface {
	Create(r *http.Request, email string, role roles.UserRole) (sqlcore2.Account, error)
	GetByEmail(r *http.Request, email string) (sqlcore2.Account, error)
	GetByID(r *http.Request, id uuid.UUID) (sqlcore2.Account, error)
	UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (sqlcore2.Account, error)
}

type accounts struct {
	queries *sqlcore2.Queries
}

func NewAccount(queries *sqlcore2.Queries) Accounts {
	return &accounts{queries: queries}
}

func (a *accounts) Create(r *http.Request, email string, role roles.UserRole) (sqlcore2.Account, error) {
	return a.queries.CreateAccount(r.Context(), sqlcore2.CreateAccountParams{
		Email: email,
		Role:  string(role),
	})
}

func (a *accounts) GetByID(r *http.Request, id uuid.UUID) (sqlcore2.Account, error) {
	return a.queries.GetAccountByID(r.Context(), id)
}

func (a *accounts) GetByEmail(r *http.Request, email string) (sqlcore2.Account, error) {
	return a.queries.GetAccountByEmail(r.Context(), email)
}

func (a *accounts) UpdateRole(r *http.Request, id uuid.UUID, role roles.UserRole) (sqlcore2.Account, error) {
	return a.queries.UpdateAccountRole(r.Context(), sqlcore2.UpdateAccountRoleParams{
		ID:   id,
		Role: string(role),
	})
}
