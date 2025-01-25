package repositories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/data/sql/repositories/sqlcore"
)

type Accounts interface {
	Create(r *http.Request, email string, role sqlcore.RoleType) (sqlcore.Account, error)

	GetById(r *http.Request, id uuid.UUID) (sqlcore.Account, error)
	GetByEmail(r *http.Request, email string) (sqlcore.Account, error)
}

type accounts struct {
	queries *sqlcore.Queries
}

func NewAccount(queries *sqlcore.Queries) Accounts {
	return &accounts{queries: queries}
}

func (a *accounts) Create(r *http.Request, email string, role sqlcore.RoleType) (sqlcore.Account, error) {
	return a.queries.CreateAccount(r.Context(), sqlcore.CreateAccountParams{
		Email: email,
		Role:  role,
	})
}

func (a *accounts) GetById(r *http.Request, id uuid.UUID) (sqlcore.Account, error) {
	return a.queries.GetAccountByID(r.Context(), id)
}

func (a *accounts) GetByEmail(r *http.Request, email string) (sqlcore.Account, error) {
	return a.queries.GetAccountByEmail(r.Context(), email)
}
